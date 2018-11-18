package cmd

// info is just a wrapper around host info, so lets 'do' that
func init() {
	hassosCmd.AddCommand(hostInfoCmd)
}
