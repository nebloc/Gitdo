const gitKey string = "Git"
	result := VCMap[gitKey].NameOfDir()
	result := VCMap[gitKey].NameOfVC()
	VCMap[gitKey].moveToDir(t)
	diff, err := VCMap[gitKey].GetDiff()
	VCMap[gitKey].moveToDir(t)
	err := VCMap[gitKey].SetHooks(HomeDir)
		filePath := filepath.Join(VCMap[gitKey].NameOfDir(), "hooks", fileName)
