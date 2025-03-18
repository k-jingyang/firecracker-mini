package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	fc "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func main() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	tempDir, err := os.MkdirTemp("", "firecracker-sock")
	if err != nil {
		log.Error().Msgf("failed to create temp dir: %v", err.Error())
	}
	fmt.Println("Temporary directory created:", tempDir)
	stdout, _ := os.Create("std-out.log")

	networkName := "ptp-net"
	config := fc.Config{
		SocketPath:      path.Join(tempDir, "firecracker.sock"),
		LogPath:         stdout.Name(),
		LogLevel:        "Info",
		KernelImagePath: "vmlinux-5.10.225",
		KernelArgs:      "console=ttyS0 reboot=k panic=1 pci=off",
		Drives: []models.Drive{
			{
				DriveID:      lo.ToPtr("rootfs"),
				PathOnHost:   lo.ToPtr("ubuntu-24.04.squashfs.upstream"),
				IsRootDevice: lo.ToPtr(true),
				IsReadOnly:   lo.ToPtr(false),
				CacheType:    lo.ToPtr("Unsafe"),
				IoEngine:     lo.ToPtr("Sync"),
				RateLimiter:  nil,
			},
		},
		MachineCfg: models.MachineConfiguration{
			VcpuCount:       lo.ToPtr(int64(2)),
			MemSizeMib:      lo.ToPtr(int64(1024)),
			Smt:             lo.ToPtr(false),
			TrackDirtyPages: false,
		},

		NetworkInterfaces: fc.NetworkInterfaces{
			fc.NetworkInterface{
				CNIConfiguration: &fc.CNIConfiguration{
					NetworkName: networkName,
					IfName:      "veth0",
				},
				AllowMMDS: true,
			},
		},
	}
	uVM, err := fc.NewMachine(context.Background(), config)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	err = uVM.Start(context.Background())

	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	// Get allocated IP address from CNI
	ipBuf, _ := os.ReadFile(fmt.Sprintf("/var/lib/cni/networks/%s/last_reserved_ip.0", networkName))
	log.Info().Msgf("IP address: %s", string(ipBuf))

	// block to let uVM spin up
	<-exitSignal
}
