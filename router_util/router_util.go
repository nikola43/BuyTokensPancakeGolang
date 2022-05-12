package genericutils

import (
	"buytokenspancakegolang/genericutils"
	"fmt"
	"time"

	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/hrharder/go-gas"
	"github.com/shopspring/decimal"
)
