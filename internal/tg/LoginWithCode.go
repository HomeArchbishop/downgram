package tg

func LoginWithCode(phone, code, phoneCodeHash string) error {
	_, err := client.Auth().SignIn(*clientCtx, phone, code, phoneCodeHash)
	if err != nil {
		return err
	}

	return nil
}
