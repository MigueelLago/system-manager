package hardware

func GetBIOSInfo() string {
	return readBiosInfo()
}
