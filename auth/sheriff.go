package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/robinstock_go"
	"github.com/yourusername/robinstock_go/urls"
)

func handleSheriffVerification(ctx context.Context, client *robinstock_go.Client, deviceToken, workflowID string) error {
	log.Println("Starting Sheriff verification workflow...")

	machinePayload := map[string]interface{}{
		"device_id": deviceToken,
		"flow":      "suv",
		"input": map[string]interface{}{
			"workflow_id": workflowID,
		},
	}

	machineResp, err := client.Post(ctx, urls.PathfinderUserMachineURL(), machinePayload, false)
	if err != nil {
		return fmt.Errorf("failed to start sheriff verification: %w", err)
	}

	if machineResp.Data == nil {
		return fmt.Errorf("no data in sheriff verification response")
	}

	machineID := robinstock_go.GetString(machineResp.Data, "id")
	if machineID == "" {
		return fmt.Errorf("no machine ID in response")
	}

	inquiryURL := urls.SheriffInquiryURL(machineID)
	log.Println("Waiting for inquiry data...")

	inquiryTimeout := time.Now().Add(20 * time.Second)
	var inquiryData map[string]interface{}

	for time.Now().Before(inquiryTimeout) {
		resp, err := client.Get(ctx, inquiryURL, nil, false)
		if err == nil && resp != nil && resp.Data != nil {
			inquiryData = resp.Data
			break
		}
		log.Println("Failed to get inquiry, retrying...")
		time.Sleep(4 * time.Second)
	}

	if inquiryData == nil {
		return fmt.Errorf("unable to get inquiry data within timeout")
	}

	contextData, ok := inquiryData["context"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no context in inquiry data")
	}

	challenge, ok := contextData["sheriff_challenge"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no sheriff_challenge in context")
	}

	challengeID := robinstock_go.GetString(challenge, "id")
	if challengeID == "" {
		return fmt.Errorf("no challenge ID")
	}

	statusURL := urls.SheriffChallengeStatusURL(challengeID)
	log.Println("Polling sheriff challenge status...")

	startTime := time.Now()
	timeout := 2 * time.Minute

	for time.Since(startTime) < timeout {
		statusResp, err := client.Get(ctx, statusURL, nil, false)
		if err != nil || statusResp == nil || statusResp.Data == nil {
			log.Println("Empty challenge status response, retrying...")
			time.Sleep(5 * time.Second)
			continue
		}

		status := robinstock_go.GetString(statusResp.Data, "challenge_status")
		log.Printf("Current Challenge Status: %s\n", status)

		switch status {
		case "validated":
			log.Println("Sheriff ID validation successful.")

			payload := map[string]interface{}{
				"sequence": 0,
				"user_input": map[string]interface{}{
					"status": "continue",
				},
			}

			finalResp, err := client.Post(ctx, inquiryURL, payload, false)
			if err != nil {
				return fmt.Errorf("failed to complete sheriff verification: %w", err)
			}

			if finalResp != nil && finalResp.Data != nil {
				typeContext, ok := finalResp.Data["type_context"].(map[string]interface{})
				if ok {
					result := robinstock_go.GetString(typeContext, "result")
					if result == "workflow_status_approved" {
						log.Println("Final workflow approval successful.")
						return nil
					}
				}
			}
			return fmt.Errorf("workflow approval failed after validation")

		case "issued":
			log.Println("Challenge still pending, waiting...")
			time.Sleep(15 * time.Second)

		default:
			return fmt.Errorf("unexpected challenge status: %s", status)
		}
	}

	return fmt.Errorf("sheriff verification timeout")
}

