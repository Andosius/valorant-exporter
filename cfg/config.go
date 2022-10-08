package cfg

import "os"

const (
	GET_PLAYER_PREFERENCES = "https://playerpreferences.riotgames.com/playerPref/v3/getPreference/Ares.PlayerSettings"
	PUT_PLAYER_PREFERENCES = "https://playerpreferences.riotgames.com/playerPref/v3/savePreference"
	CLIENT_PLATFORM        = "ew0KCSJwbGF0Zm9ybVR5cGUiOiAiUEMiLA0KCSJwbGF0Zm9ybU9TIjogIldpbmRvd3MiLA0KCSJwbGF0Zm9ybU9TVmVyc2lvbiI6ICIxMC4wLjE5MDQyLjEuMjU2LjY0Yml0IiwNCgkicGxhdGZvcm1DaGlwc2V0IjogIlVua25vd24iDQp9"

	DATA_DIR = "data"
	PERMS    = 0660
	SEP      = string(os.PathSeparator)
)
