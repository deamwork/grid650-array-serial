package runtime

type Flag struct {
	HTTPListen  string      `flag:"name=http-listen" flag-usage:"Specific http listen address."`
	Device      string      `flag:"env name=device" flag-usage:"Specific device tty"`
	Baud        int         `flag:"env name=baud" flag-usage:"Specific device bit rate"`
	ConfigFile  string      `flag:"name=config" flag-usage:"Specific config file path."`
	Environment Environment `flag:""`
}

type Environment struct {
	DevelopmentMode bool `flag:"env short=d" flag-usage:"change environment mode to development"`
	JSONLogStyle    bool `flag:"env name=json-log-style short=j" flag-usage:"change log style to JSON"`
}
