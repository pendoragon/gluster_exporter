package main

import (
	"bytes"
	"os/exec"
	"strconv"

	"github.com/ofesseler/gluster_exporter/structs"
	"github.com/prometheus/common/log"
)

func execGlusterCommand(arg ...string) (*bytes.Buffer, error) {
	stdoutBuffer := &bytes.Buffer{}
	argXML := append(arg, "--xml")
	glusterExec := exec.Command(GlusterCmd, argXML...)
	glusterExec.Stdout = stdoutBuffer
	err := glusterExec.Run()

	if err != nil {
		log.Errorf("tried to execute %v and got error: %v", arg, err)
		return stdoutBuffer, err
	}
	return stdoutBuffer, nil
}

// ExecVolumeInfo executes "gluster volume info" at the local machine and
// returns VolumeInfoXML struct and error
func ExecVolumeInfo() (structs.VolumeInfoXML, error) {
	args := []string{"volume", "info"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return structs.VolumeInfoXML{}, cmdErr
	}
	volumeInfo, err := structs.VolumeInfoXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeInfo, err
	}

	return volumeInfo, nil
}

// ExecVolumeList executes "gluster volume info" at the local machine and
// returns VolumeList struct and error
func ExecVolumeList() (structs.VolList, error) {
	args := []string{"volume", "list"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return structs.VolList{}, cmdErr
	}
	volumeList, err := structs.VolumeListXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeList.VolList, err
	}

	return volumeList.VolList, nil
}

// ExecPeerStatus executes "gluster peer status" at the local machine and
// returns PeerStatus struct and error
func ExecPeerStatus() (structs.PeerStatus, error) {
	args := []string{"peer", "status"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return structs.PeerStatus{}, cmdErr
	}
	peerStatus, err := structs.PeerStatusXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return peerStatus.PeerStatus, err
	}

	return peerStatus.PeerStatus, nil
}

// ExecVolumeProfileGvInfoCumulative executes "gluster volume {volume] profile info cumulative" at the local machine and
// returns VolumeInfoXML struct and error
func ExecVolumeProfileGvInfoCumulative(volumeName string) (structs.VolProfile, error) {
	args := []string{"volume", "profile"}
	args = append(args, volumeName)
	args = append(args, "info", "cumulative")
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return structs.VolProfile{}, cmdErr
	}
	volumeProfile, err := structs.VolumeProfileGvInfoCumulativeXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeProfile.VolProfile, err
	}
	return volumeProfile.VolProfile, nil
}

// ExecVolumeStatusAllDetail executes "gluster volume status all detail" at the local machine
// returns VolumeStatusXML struct and error
func ExecVolumeStatusAllDetail() (structs.VolumeStatusXML, error) {
	args := []string{"volume", "status", "all", "detail"}
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return structs.VolumeStatusXML{}, cmdErr
	}
	volumeStatus, err := structs.VolumeStatusAllDetailXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Errorf("Something went wrong while unmarshalling xml: %v", err)
		return volumeStatus, err
	}
	return volumeStatus, nil
}

// ExecVolumeHealInfo executes volume heal info on host system and processes input
// returns (int) number of unsynced files
func ExecVolumeHealInfo(volumeName string) (int, error) {
	args := []string{"volume", "heal", volumeName, "info"}
	entriesOutOfSync := 0
	bytesBuffer, cmdErr := execGlusterCommand(args...)
	if cmdErr != nil {
		return -1, cmdErr
	}
	healInfo, err := structs.VolumeHealInfoXMLUnmarshall(bytesBuffer)
	if err != nil {
		log.Error(err)
		return -1, err
	}

	for _, brick := range healInfo.HealInfo.Bricks.Brick {
		var count int
		count, _ = strconv.Atoi(brick.NumberOfEntries)
		entriesOutOfSync += count
	}
	return entriesOutOfSync, nil
}
