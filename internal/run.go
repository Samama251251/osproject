package internal

import (
	"fmt"
	"github.com/samama/firaaq/pkg/container"
	"github.com/samama/firaaq/pkg/image"
	"github.com/samama/firaaq/pkg/network"
	"github.com/samama/firaaq/pkg/reexec"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"syscall"
)

// Run: runs a command inside a new container.
func Run(cmd *cobra.Command, args []string) error {
	ctr := container.NewContainer()
	defer ctr.Remove()

	// setup bridge
	if err := network.SetupBridge("firaaq0"); err != nil {
		return err
	}

	// setup network
	deleteNetwork, err := ctr.SetupNetwork("firaaq0")
	if err != nil {
		return err
	}
	defer deleteNetwork()

	// setup filesystem from image
	img, err := getImage(args[0])
	if err != nil {
		return err
	}
	unmount, err := ctr.MountFromImage(img)
	if err != nil {
		return errors.Wrap(err, "can't Mount image filesystem")
	}
	defer unmount()

	// Format fork options so the child process sees the same runtime flags plus mount info.
	options := append([]string{}, rawFlags(cmd.Flags())...)
	options = append(options, fmt.Sprintf("--root=%s", ctr.RootFS))
	options = append(options, fmt.Sprintf("--container=%s", ctr.Digest))
	newArgs := []string{"fork"}
	newArgs = append(newArgs, options...)
	newArgs = append(newArgs, args[1:]...)
	newCmd := reexec.Command(newArgs...)
	// Re-execute ourselves inside the fork helper so the namespaces are applied correctly.
	newCmd.Stdin, newCmd.Stdout, newCmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	newCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID,
	}
	if err := newCmd.Run(); err != nil {
		return errors.Wrap(err, "failed run fork process")
	}

	return nil
}

// rawFlags convert a pflag.FlagSet to a slice of string.
func rawFlags(flags *pflag.FlagSet) []string {
	var flagList []string
	flags.VisitAll(func(flag *pflag.Flag) {
		if flag.Value.String() == "" {
			return
		}
		flagList = append(flagList, fmt.Sprintf("--%s=%v", flag.Name, flag.Value))
	})
	return flagList
}

// getImage: pulls an image and download its layers if it image does not exist locally.
func getImage(src string) (*image.Image, error) {
	img, err := image.NewImage(src)
	if err != nil {
		return img, errors.Wrapf(err, "Can't pull %q", src)
	}
	exists, err := img.Exists()
	if err != nil {
		return img, err
	}
	if !exists {
		fmt.Printf("Unable to find image %s:%s locally\n", img.Repository, img.Tag)
		fmt.Printf("downloading the image from %s\n", img.Registry)
		err = img.Download()
	}

	return img, err
}

