package handlers

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dbyington/httpio"
	"github.com/dbyington/manifestgo"
	"github.com/sirupsen/logrus"

	"github.com/dbyington/manifest-server/store"
	"github.com/dbyington/manifest-server/store/redis"
)

type PkgHandler struct {
	store.Store
	*logrus.Logger
}

func NewPkgHandler() *PkgHandler {
	return &PkgHandler{
		redis.New(),
		logrus.New(),
	}
}

func (h *PkgHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		o   *store.Object
	)

	h.Info("handling request")

	// Kludge
	param, ok := r.URL.Query()["pkgurl"]
	if !ok {
		h.Error("while getting package url")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pkgUrl := param[0]

	hashParam, ok := r.URL.Query()["hashtype"]
	if !ok {
		h.Error("while getting hash type")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var (
		hashSize uint
		hashType string
	)
	switch hashParam[0] {
	case "sha256":
		hashSize = sha256.Size
		hashType = "sha256"
	default:
		hashSize = md5.Size
		hashType = "md5"
	}

	// Since having multiple hash types creates multiple assets, which seems to mess with installing, create one per hash type.
	cacheKey := fmt.Sprintf("%s+%d", pkgUrl, hashSize)

	expectHeaders := make(map[string]string)

	h.Info("checking cache...")
	o, err = h.Get(r.Context(), cacheKey)
	if err != nil {
		if err != store.ErrCacheMiss {
			h.Errorf("while checking cache: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.Infof("missed cache: %s, must read fresh", err)
	} else {
	    h.Infof("GOT OBJECT WITH ETAG: %s", o.Etag)
	    // We got a hit in the cache, so check against the Etag.
		expectHeaders = map[string]string{
			"Etag": o.Etag,
			// for now just try with etag
			// "ContentLength": strconv.FormatInt(o.ContentLength, 10),
		}
	}

	h.Info("reading package url...")
	p, err := manifestgo.ReadPkgUrl(&http.Client{}, pkgUrl, hashSize, expectHeaders)
	if err != nil {
		if err != httpio.ErrHeaderEtag {
			h.Errorf("while reading package: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// We only need to do this if the Etag didn't match.
		h.Info("building manifest...")
		m, err := p.BuildManifest()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.Info("creating json package....")
		asJson, err := m.AsJSON(0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.Info("creating plit package...")
		asPlist, err := m.AsEncodedPlistString(4)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.Infof("P ETAG: %s", p.Etag)

		// If we got a new response let's update our o instance.
		o = &store.Object{
			ContentLength: p.ContentLength,
			Etag:          p.Etag,
			HashType:      hashType,
			Id:            p.Etag + strconv.FormatUint(uint64(hashSize), 10),
			Json:          asJson,
			Plist:         asPlist,
			Title:         p.Title,
		}
	}

	h.Infof("O ETAG: %s", o.Etag)

	// Putting it back here will restart the TTL, keeping it in cache longer.
	h.Info("putting back into cache")
	if err = h.Put(r.Context(), cacheKey, o); err != nil {
		h.Errorf("while storing manifest: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		// , "error while caching manifest")
		return
	}

	h.Info("ALMOST DONE, preparing response...")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	if err := e.Encode(o); err != nil {
		h.Errorf("while encoding response: %s", err)
	}
}
