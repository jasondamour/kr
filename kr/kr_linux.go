package main

import (
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

func runCommandWithUserInteraction(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func restartCommand(c *cli.Context) (err error) {
	exec.Command("systemctl", "--user", "daemon-reload").Run()
	exec.Command("systemctl", "--user", "disable", "kr").Run()
	exec.Command("systemctl", "--user", "stop", "kr").Run()
	exec.Command("systemctl", "--user", "enable", "kr").Run()
	if err := exec.Command("systemctl", "--user", "start", "kr").Run(); err != nil {
		//	fall back to system-level daemon
		runCommandWithUserInteraction("systemctl", "restart", "kr")
	}
	PrintErr(os.Stderr, "Restarted Kryptonite daemon.")
	return
}

func openBrowser(url string) {
	exec.Command("sensible-browser", url).Run()
}

func uninstallCommand(c *cli.Context) (err error) {
	confirmOrFatal(os.Stderr, "Uninstall Kryptonite from this workstation? (same as sudo apt-get/yum remove kr)")

	exec.Command("systemctl", "--user", "disable", "kr").Run()
	if err := exec.Command("systemctl", "--user", "stop", "kr").Run(); err != nil {
		runCommandWithUserInteraction("systemctl", "stop", "kr")
	}

	if aptGetErr := exec.Command("which", "apt-get").Run(); aptGetErr == nil {
		uninstallCmd := exec.Command("sudo", "apt-get", "remove", "kr", "-y")
		uninstallCmd.Stdout = os.Stdout
		uninstallCmd.Stderr = os.Stderr
		uninstallCmd.Run()
	}

	if yumErr := exec.Command("which", "yum").Run(); yumErr == nil {
		uninstallCmd := exec.Command("sudo", "yum", "remove", "kr", "-y")
		uninstallCmd.Stdout = os.Stdout
		uninstallCmd.Stderr = os.Stderr
		uninstallCmd.Run()
	}

	PrintErr(os.Stderr, "Kryptonite uninstalled.")
	return
}

func upgradeCommand(c *cli.Context) (err error) {
	confirmOrFatal(os.Stderr, "Upgrade Kryptonite on this workstation?")
	update := exec.Command("sudo", "apt-get", "update")
	update.Stdout = os.Stdout
	update.Stderr = os.Stderr
	update.Stdin = os.Stdin
	update.Run()
	cmd := exec.Command("sudo", "apt-get", "install", "kr")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
	return
}
