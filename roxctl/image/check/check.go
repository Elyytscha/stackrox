package check

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/images/utils"
	"github.com/stackrox/rox/pkg/retry"
	pkgUtils "github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/roxctl/common"
	"github.com/stackrox/rox/roxctl/common/flags"
	"github.com/stackrox/rox/roxctl/common/report"
	"github.com/stackrox/rox/roxctl/common/util"
)

// Command checks the image against image build lifecycle policies
func Command() *cobra.Command {
	var (
		image             string
		json              bool
		retryDelay        int
		retryCount        int
		sendNotifications bool
	)
	c := &cobra.Command{
		Use: "check",
		RunE: util.RunENoArgs(func(c *cobra.Command) error {
			return checkImageWithRetry(image, json, sendNotifications, flags.Timeout(c), retryDelay, retryCount)
		}),
	}

	c.Flags().StringVarP(&image, "image", "i", "", "image name and reference. (e.g. nginx:latest or nginx@sha256:...)")
	pkgUtils.Must(c.MarkFlagRequired("image"))

	c.Flags().BoolVar(&json, "json", false, "output policy results as json.")
	c.Flags().IntVarP(&retryDelay, "retry-delay", "d", 3, "set time to wait between retries in seconds")
	c.Flags().IntVarP(&retryCount, "retries", "r", 0, "Number of retries before exiting as error")
	c.Flags().BoolVar(&sendNotifications, "send-notifications", false,
		"whether to send notifications for violations (notifications will be sent to the notifiers "+
			"configured in each violated policy)")
	return c
}

func checkImageWithRetry(image string, json bool, sendNotifications bool, timeout time.Duration, retryDelay int, retryCount int) error {
	err := retry.WithRetry(func() error {
		return checkImage(image, json, sendNotifications, timeout)
	},
		retry.Tries(retryCount+1),
		retry.OnFailedAttempts(func(err error) {
			fmt.Fprintf(os.Stderr, "Checking image failed: %v. Retrying after %v seconds\n", err, retryDelay)
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}))
	if err != nil {
		return err
	}
	return nil
}

func checkImage(image string, json bool, sendNotifications bool, timeout time.Duration) error {
	// Get the violated policies for the input data.
	req, err := buildRequest(image, sendNotifications)
	if err != nil {
		return err
	}
	alerts, err := sendRequestAndGetAlerts(req, timeout)
	if err != nil {
		return err
	}

	// If json mode was given, print results (as json) and immediately return.
	if json {
		return report.JSON(os.Stdout, alerts)
	}

	// Print results in human readable mode.
	if err = report.PrettyWithResourceName(os.Stdout, alerts, storage.EnforcementAction_FAIL_BUILD_ENFORCEMENT, "Image", image); err != nil {
		return err
	}

	// Check if any of the violated policies have an enforcement action that
	// fails the CI build.
	for _, alert := range alerts {
		if report.EnforcementFailedBuild(storage.EnforcementAction_FAIL_BUILD_ENFORCEMENT)(alert.GetPolicy()) {
			return errors.New("Violated a policy with CI enforcement set")
		}
	}
	return nil
}

// Get the alerts for the command line inputs.
func sendRequestAndGetAlerts(req *v1.BuildDetectionRequest, timeout time.Duration) ([]*storage.Alert, error) {
	// Create the connection to the central detection service.
	conn, err := common.GetGRPCConnection()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()
	service := v1.NewDetectionServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Call detection and return the returned alerts.
	response, err := service.DetectBuildTime(ctx, req)
	if err != nil {
		return nil, err
	}
	return response.GetAlerts(), nil
}

// Use inputs to generate an image name for request.
func buildRequest(image string, sendNotifications bool) (*v1.BuildDetectionRequest, error) {
	img, err := utils.GenerateImageFromString(image)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse image '%s'", image)
	}
	return &v1.BuildDetectionRequest{
		Resource:          &v1.BuildDetectionRequest_Image{Image: img},
		SendNotifications: sendNotifications,
	}, nil
}
