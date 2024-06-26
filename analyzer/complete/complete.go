package complete

func Check(resultFolderPath string, progPath string) {
	progElems, err := getProgramElements(progPath)
	if err != nil {
		panic(err)
	}

	traceElems, err := getTraceElements(resultFolderPath)
	if err != nil {
		panic(err)
	}

	for file, lines := range traceElems {
		print(file, " -> ", lines)
	}

	_ = progElems
	_ = traceElems
}
