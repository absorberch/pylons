package keep

import (
	"strings"
	"testing"

	"github.com/MikeSofaer/pylons/x/pylons/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func GenExecution(sender sdk.AccAddress, tci TestCoinInput) types.Execution {
	cbData := GenCookbook(sender, "cookbook-0001", "this has to meet character limits")
	rcpData := GenRecipe(sender, cbData.ID, "new recipe", "this has to meet character limits; lol")

	var cl, ocl sdk.Coins
	for _, inp := range rcpData.CoinInputs {
		cl = append(cl, sdk.NewCoin(inp.Coin, sdk.NewInt(inp.Count)))
	}
	for _, out := range rcpData.CoinOutputs {
		ocl = append(ocl, sdk.NewCoin(out.Coin, sdk.NewInt(out.Count)))
	}

	var inputItems []types.Item

	inputItems = append(inputItems, *GenItem(cbData.ID, sender, "Raichu"))
	inputItems = append(inputItems, *GenItem(cbData.ID, sender, "Raichu"))

	var outputItems []types.Item

	outputItems = append(inputItems, *GenItem(cbData.ID, sender, "Zombie"))

	exec := types.Execution{
		RecipeID:    rcpData.ID,
		CoinInputs:  cl,
		CoinOutputs: ocl,
		BlockHeight: tci.Ctx.BlockHeight() + rcpData.BlockInterval,
		ItemInputs:  inputItems,
		ItemOutputs: outputItems,
		Sender:      sender,
		Completed:   false,
	}
	exec.ID = exec.KeyGen()
	return exec
}

func TestKeeperSetExecution(t *testing.T) {
	mockedCoinInput := SetupTestCoinInput()

	sender, _ := sdk.AccAddressFromBech32("cosmos1y8vysg9hmvavkdxpvccv2ve3nssv5avm0kt337")

	cases := map[string]struct {
		sender       sdk.AccAddress
		desiredError string
		showError    bool
	}{
		"empty sender test": {
			sender:       nil,
			desiredError: "SetExecution: the sender cannot be empty",
			showError:    true,
		},
		"successful set execution test": {
			sender:       sender,
			desiredError: "",
			showError:    false,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			exec := GenExecution(tc.sender, mockedCoinInput)
			err := mockedCoinInput.PlnK.SetExecution(mockedCoinInput.Ctx, exec)

			if tc.showError {
				// t.Errorf("execution_test err LOG:: %+v", err)
				require.True(t, strings.Contains(err.Error(), tc.desiredError))
			} else {
				require.True(t, err == nil)
			}
		})
	}
}

func TestKeeperGetExecution(t *testing.T) {
	mockedCoinInput := SetupTestCoinInput()

	sender, _ := sdk.AccAddressFromBech32("cosmos1y8vysg9hmvavkdxpvccv2ve3nssv5avm0kt337")
	exec := GenExecution(sender, mockedCoinInput)
	mockedCoinInput.PlnK.SetExecution(mockedCoinInput.Ctx, exec)

	cases := map[string]struct {
		execId       string
		desiredError string
		showError    bool
	}{
		"wrong exec id test": {
			execId:       "invalidExecID",
			desiredError: "The execution doesn't exist",
			showError:    true,
		},
		"successful get execution test": {
			execId:       exec.ID,
			desiredError: "",
			showError:    false,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			execution, err := mockedCoinInput.PlnK.GetExecution(mockedCoinInput.Ctx, tc.execId)

			if tc.showError {
				require.True(t, strings.Contains(err.Error(), tc.desiredError))
			} else {
				// t.Errorf("execution_test err LOG:: %+v", err)
				require.True(t, err == nil)
				require.True(t, execution.RecipeID == exec.RecipeID)
				require.True(t, execution.Sender.String() == exec.Sender.String())
				require.True(t, execution.Completed == exec.Completed)
			}
		})
	}
}

func TestKeeperUpdateExecution(t *testing.T) {
	mockedCoinInput := SetupTestCoinInput()

	sender, _ := sdk.AccAddressFromBech32("cosmos1y8vysg9hmvavkdxpvccv2ve3nssv5avm0kt337")
	exec := GenExecution(sender, mockedCoinInput)
	mockedCoinInput.PlnK.SetExecution(mockedCoinInput.Ctx, exec)
	newExec := GenExecution(sender, mockedCoinInput)
	newExec.Completed = true

	cases := map[string]struct {
		execId       string
		desiredError string
		showError    bool
	}{
		"wrong exec id test": {
			execId:       "invalidExecID",
			desiredError: "the exec with gid invalidExecID does not exist",
			showError:    true,
		},
		"successful update execution test": {
			execId:       exec.ID,
			desiredError: "",
			showError:    false,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			err := mockedCoinInput.PlnK.UpdateExecution(mockedCoinInput.Ctx, tc.execId, newExec)

			if tc.showError {
				require.True(t, strings.Contains(err.Error(), tc.desiredError))
			} else {
				// t.Errorf("execution_test err LOG:: %+v", err)
				require.True(t, err == nil)
				uExec, err2 := mockedCoinInput.PlnK.GetExecution(mockedCoinInput.Ctx, tc.execId)
				require.True(t, err2 == nil)
				require.True(t, uExec.Completed == true)
			}
		})
	}
}
