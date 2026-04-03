package timeentry

func initCreateFlags() {
	createCmd.Flags().Uint64VarP(&createWorkPackageId, "work-package", "w", 0, "Work package ID to log time on")
	createCmd.Flags().Float64VarP(&createHours, "hours", "", 0, "Number of hours spent (e.g. 1.5 for 1h30m)")
	createCmd.Flags().StringVarP(&createActivity, "activity", "", "", "Activity name (e.g. Development, Design)")
	createCmd.Flags().StringVarP(&createSpentOn, "spent-on", "", "", "Date of the time entry (YYYY-MM-DD, default: today)")
	createCmd.Flags().Uint64VarP(&createUserId, "user", "u", 0, "User ID (default: current user)")
	createCmd.Flags().StringVarP(&createComment, "comment", "", "", "Comment for the time entry")

	_ = createCmd.MarkFlagRequired("work-package")
	_ = createCmd.MarkFlagRequired("hours")
}
