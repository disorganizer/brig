package testlayer

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/disorganizer/brig/repo"
	"github.com/disorganizer/brig/transfer"
	"github.com/disorganizer/brig/transfer/moose"
	"github.com/disorganizer/brig/util/testwith"
)

func WithConnector(t *testing.T, user string, fc func(c *transfer.Connector)) {
	testwith.WithRepo(t, user, user+"pass", func(rp *repo.Repository) {
		if err := rp.IPFS.Online(); err != nil {
			t.Errorf("Cannot go online with IPFS repo: %v", err)
			return
		}

		lay := moose.NewLayer(rp.IPFS, context.Background())
		con := transfer.NewConnector(lay, rp)

		if err := con.Connect(); err != nil {
			t.Errorf("Cannot connect: %v", err)
			return
		}

		t.Logf("Entering test for %s's connector", user)
		fc(con)
		t.Logf("Leaving test for %s's connector", user)

		if err := con.Disconnect(); err != nil {
			t.Errorf("Cannot disconnect: %v", err)
			return
		}
	})
}

func TestConversation(t *testing.T) {
	WithConnector(t, "alice", func(ac *transfer.Connector) {
		WithConnector(t, "bob", func(bc *transfer.Connector) {
			br, ar := bc.Repo(), ac.Repo()
			berr := br.Remotes.Insert(repo.NewRemoteFromPeer(ar.Peer()))
			if berr != nil {
				t.Errorf("Bob has no friends: %v", berr)
				return
			}

			aerr := ar.Remotes.Insert(repo.NewRemoteFromPeer(br.Peer()))
			if aerr != nil {
				t.Errorf("Alice has no friends: %v", aerr)
				return
			}

			apc, err := ac.Dial(br.Peer())
			if err != nil {
				t.Errorf("Alice cannot dial to bob: %v", err)
				return
			}

			// Spam in some queries:
			for i := 0; i < 10; i++ {
				v, err := apc.QueryStoreVersion()
				if err != nil {
					t.Errorf("Usage of api client failed: %v", err)
					return
				}

				if v <= 0 {
					t.Errorf("Version should be any positive number")
					return
				}
			}

			if err := apc.Close(); err != nil {
				t.Errorf("Alice cannot close apiclient to bob: %v", err)
				return
			}
		})
	})
}
