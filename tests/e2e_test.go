package tests

import (
	"github.com/alexisvisco/graphql-proxy-grpc/internal/gengql"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestSimpleSuperheroService(t *testing.T) {
	err := os.Chdir("simple-superhero-service")
	assert.NoError(t, err)

	cmd(t, "rm -rf gen")
	cmd(t, "buf generate --template buf.gen.yaml")

	gengql.GenerateGoGqlFromDir("./gen/go/graphql")

	checkSnapshot(t, "graphql-schema.graphql", "./gen/go/graphql/schema.graphql")
	checkSnapshot(t, "graphql-models_gen.go", "./gen/go/graphql/models_gen.go")
	checkSnapshot(t, "graphql-resolvers.go", "./gen/go/graphql/resolvers.go")
	checkSnapshot(t, "graphql-schema.go", "./gen/go/graphql/schema.go")

	checkSnapshot(t, "protos-enums.go", "./gen/go/protos/superheroes/enums.gql.go")
}

func checkSnapshot(t *testing.T, snapshotID string, file string) {
	t.Run(snapshotID, func(t *testing.T) {
		content, err := ioutil.ReadFile(file)
		assert.NoError(t, err)

		err = cupaloy.SnapshotMulti(snapshotID, content)
		assert.NoError(t, err)
	})
}

func cmd(t *testing.T, command string) {
	line := strings.Split(command, " ")
	cmd := exec.Command(line[0], line[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	assert.NoError(t, err)
}
