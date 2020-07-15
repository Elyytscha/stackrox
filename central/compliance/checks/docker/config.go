package docker

import (
	"github.com/stackrox/rox/central/compliance/checks/common"
	"github.com/stackrox/rox/central/compliance/framework"
)

func init() {
	framework.MustRegisterChecksIfFlagDisabled(
		common.SystemdOwnershipCheck("CIS_Docker_v1_2_0:3_1", "docker.service", "root", "root"),
		common.SystemdPermissionCheck("CIS_Docker_v1_2_0:3_2", "docker.service", 0644),

		common.OptionalSystemdOwnershipCheck("CIS_Docker_v1_2_0:3_3", "docker.socket", "root", "root"),
		common.OptionalSystemdPermissionCheck("CIS_Docker_v1_2_0:3_4", "docker.socket", 0644),

		common.OwnershipCheck("CIS_Docker_v1_2_0:3_5", "/etc/docker", "root", "root"),
		common.PermissionCheck("CIS_Docker_v1_2_0:3_6", "/etc/docker", 0755),

		common.RecursiveOwnershipCheck("CIS_Docker_v1_2_0:3_7", "/etc/docker/certs.d", "root", "root"),
		common.RecursivePermissionCheck("CIS_Docker_v1_2_0:3_8", "/etc/docker/certs.d", 0444),

		common.CommandLineFileOwnership("CIS_Docker_v1_2_0:3_9", "dockerd", "tlscacert", "root", "root"),
		common.CommandLineFilePermissions("CIS_Docker_v1_2_0:3_10", "dockerd", "tlscacert", 0444),

		common.CommandLineFileOwnership("CIS_Docker_v1_2_0:3_11", "dockerd", "tlscert", "root", "root"),
		common.CommandLineFilePermissions("CIS_Docker_v1_2_0:3_12", "dockerd", "tlscert", 0444),

		common.CommandLineFileOwnership("CIS_Docker_v1_2_0:3_13", "dockerd", "tlskey", "root", "root"),
		common.CommandLineFilePermissions("CIS_Docker_v1_2_0:3_14", "dockerd", "tlskey", 0400),

		common.OwnershipCheck("CIS_Docker_v1_2_0:3_15", "/var/run/docker.sock", "root", "docker"),
		common.PermissionCheck("CIS_Docker_v1_2_0:3_16", "/var/run/docker.sock", 0660),

		common.OptionalOwnershipCheck("CIS_Docker_v1_2_0:3_17", "/etc/docker/daemon.json", "root", "root"),
		common.OptionalPermissionCheck("CIS_Docker_v1_2_0:3_18", "/etc/docker/daemon.json", 0644),

		common.OptionalOwnershipCheck("CIS_Docker_v1_2_0:3_19", "/etc/default/docker", "root", "root"),

		common.OptionalOwnershipCheck("CIS_Docker_v1_2_0:3_20", "/etc/sysconfig/docker", "root", "root"),
		common.OptionalPermissionCheck("CIS_Docker_v1_2_0:3_21", "/etc/sysconfig/docker", 0644),

		common.OptionalPermissionCheck("CIS_Docker_v1_2_0:3_22", "/etc/default/docker", 0644),
	)
}
