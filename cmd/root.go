/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/jeremyrickard/notation-attach/pkg/notation"
	"github.com/jeremyrickard/notation-attach/pkg/oras"
	"github.com/spf13/cobra"
)

type pushSignatureCmd struct {
	artifact  string
	signature string
	notation  notation.Interface
	oras      oras.Interface
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	p := pushSignatureCmd{}

	cmd := &cobra.Command{
		Use:     "notation-attach",
		Short:   "A tool for attaching notation signatures to an artifact",
		Long:    "A tool for attaching notation signatures to an artifact",
		PreRunE: p.validate,
		RunE:    p.run,
	}
	f := cmd.Flags()
	f.StringVarP(&p.artifact, "artifact", "a", "", "artifact reference")
	f.StringVarP(&p.signature, "signature", "s", "", "artifact subject reference")

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func (p *pushSignatureCmd) validate(cmd *cobra.Command, args []string) error {
	return nil
}

func (p *pushSignatureCmd) run(cmd *cobra.Command, args []string) error {
	return nil
}
