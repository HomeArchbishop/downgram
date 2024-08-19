package tg

func GetStatus() (bool, error) {
	status, err := client.Auth().Status(*clientCtx)
	if err != nil {
		return false, err
	}
	if !status.Authorized {
		return false, nil
	}
	return status.Authorized, nil
}
