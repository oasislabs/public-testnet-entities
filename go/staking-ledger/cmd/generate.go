package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/oasislabs/oasis-core/go/common/logging"
	nodeCmdCommon "github.com/oasislabs/oasis-core/go/oasis-node/cmd/common"
	"github.com/oasislabs/the-quest-entities/go/staking-ledger/stakinggenesis"
)

const (
	cfgFaucetAddress           = "faucet.address"
	cfgFaucetAmount            = "faucet.amount"
	cfgTotalSupply             = "total-supply"
	cfgPrecisionConstant       = "precision-constant"
	cfgEntitiesDirectoryPath   = "entities-dir-path"
	cfgConsensusParametersPath = "consensus-params-path"
	cfgDefaultFundingAmount    = "default-funding"
	cfgDefaultSelfEscrowAmount = "default-self-escrow"
	cfgOutputPath              = "output-path"
	defaultPrecisionConstant   = 1_000_000_000_000_000_000
	defaultTotalSupply         = 10_000_000_000
	defaultSelfEscrowAmount    = 100
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates a staking ledger",
		Long: `Generates a staking ledger

		Uses a directory of unpacked Entity Packages.
		Amounts are configured in whole tokens`,
		Run: doGenerate,
	}

	generateFlags = flag.NewFlagSet("", flag.ContinueOnError)
	logger        = logging.GetLogger("cmd/staking-ledger")
)

func doGenerate(cmd *cobra.Command, args []string) {
	options := stakinggenesis.GenesisOptions{
		FaucetBase64Address:     viper.GetString(cfgFaucetAddress),
		FaucetAmount:            viper.GetInt64(cfgFaucetAmount),
		TotalSupply:             viper.GetInt64(cfgTotalSupply),
		PrecisionConstant:       viper.GetInt64(cfgPrecisionConstant),
		EntitiesDirectoryPath:   viper.GetString(cfgEntitiesDirectoryPath),
		ConsensusParametersPath: viper.GetString(cfgConsensusParametersPath),
		DefaultFundingAmount:    viper.GetInt64(cfgDefaultFundingAmount),
		DefaultSelfEscrowAmount: viper.GetInt64(cfgDefaultSelfEscrowAmount),
	}

	if err := nodeCmdCommon.Init(); err != nil {
		nodeCmdCommon.EarlyLogAndExit(err)
	}

	outputPath := viper.GetString(cfgOutputPath)
	if outputPath == "" {
		logger.Error("must set output path for staking genesis file")
		os.Exit(1)
	}

	if options.EntitiesDirectoryPath == "" {
		logger.Error("must define an entities directory path")
		os.Exit(1)
	}
	entitiesDir, err := stakinggenesis.LoadEntitiesDirectory(options.EntitiesDirectoryPath)
	if err != nil {
		logger.Error("Cannot load entities",
			"err", err,
		)
		os.Exit(1)
	}
	options.Entities = entitiesDir

	stakingGenesis, err := stakinggenesis.Create(options)
	if err != nil {
		logger.Error("failed to create a staking genesis file",
			"err", err,
		)
		os.Exit(1)
	}

	b, err := json.Marshal(stakingGenesis)
	err = ioutil.WriteFile(outputPath, b, 0644)
	if err != nil {
		logger.Error("failed to write staking genesis to json",
			"err", err,
		)
		os.Exit(1)
	}
}

// RegisterForTestingCmd registers the for-testing subcommand.
func RegisterGenerateCmd(parentCmd *cobra.Command) {
	generateFlags.Int64(cfgFaucetAmount, 0, "amount to fund (in whole tokens)")
	generateFlags.String(cfgFaucetAddress, "", "faucet address (base64 encoded)")
	generateFlags.Int64(cfgTotalSupply, defaultTotalSupply, "Total supply of tokens (in whole tokens)")
	generateFlags.Int64(cfgPrecisionConstant, defaultPrecisionConstant,
		"the precision constant for a single token defaults to 10^18")
	generateFlags.String(cfgEntitiesDirectoryPath, "", "a directory entities")
	generateFlags.String(cfgConsensusParametersPath, "",
		"a consensus params json file (defaults to using ./consensus_params.json relative to entities directory)")
	generateFlags.Int64(cfgDefaultFundingAmount, 0, "Default funding amount")
	generateFlags.Int64(cfgDefaultSelfEscrowAmount, defaultSelfEscrowAmount, "Default amount to self escrow")
	generateFlags.String(cfgOutputPath, "", "output path for the staking ledger")
	_ = viper.BindPFlags(generateFlags)

	generateCmd.Flags().AddFlagSet(generateFlags)

	parentCmd.AddCommand(generateCmd)
}
