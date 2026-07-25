package timeentry

func initCreateFlags() {
	createCmd.Flags().StringVarP(&createWorkPackageId, "work-package", "w", "", "Work package ID or identifier to log time on")
	createCmd.Flags().Float64VarP(&createHours, "hours", "", 0, "Number of hours spent (e.g. 1.5 for 1h30m)")
	createCmd.Flags().StringVarP(&createActivity, "activity", "", "", "Activity name (e.g. Development, Design)")
	createCmd.Flags().StringVarP(&createSpentOn, "spent-on", "", "", "Date of the time entry (YYYY-MM-DD, default: today)")
	createCmd.Flags().Uint64VarP(&createUserId, "user", "u", 0, "User ID (default: current user)")
	createCmd.Flags().StringVarP(&createComment, "comment", "", "", "Comment for the time entry")

	_ = createCmd.MarkFlagRequired("work-package")
	_ = createCmd.MarkFlagRequired("hours")
}
