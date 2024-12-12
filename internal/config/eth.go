package config

type EthConfig struct {
	URL string `env:"ETH_URL" envDefault:"wss://sepolia.infura.io/ws/v3/6c89fb7fa351451f939eea9da6bee755"`
}
