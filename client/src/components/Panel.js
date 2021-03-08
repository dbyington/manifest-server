import { useState, useEffect } from 'react';

import './Panel.css';
import History from "./History";
import Current from "./Current";

const _ = require('lodash')

const serverAddr = (process.env.NODE_ENV === "development") ? 'http://localhost:8080/manifest' : window.origin + "/manifest"

export default function Panel() {
  const key = 'manifests'
  const [manifests, setManifests] = useState(() => {
    const localCache = window.localStorage.getItem(key);
    return localCache !== null ? JSON.parse(localCache) : [];
  });

  const [manifest, setManifest] = useState(null);
  const [loading, setLoading] = useState(false)

  // Load items from local storage.
  useEffect(() => {
    if (manifests.length > 0) window.localStorage.setItem(key, JSON.stringify(manifests));
  }, [key, manifests]);

  const showManifest = (m) => setManifest(m)

  // Used to retrieve manifest from the server.
  const getManifest = async (data) => {
    setLoading(true)
    setManifest({})
    console.log(`GET ${data.pkgUrl} with hash type of ${data.hashType}`)

    const url = new URL(serverAddr)
    const params = {
      pkgurl: data.pkgUrl,
      hashtype: data.hashType,
    }
    url.search = new URLSearchParams(params).toString();
    console.log(`Make Request: ${url.toString()}`)
    fetch(url.toString(), {
      method: 'GET',
      headers: {"accept": "application/json", 'content-type': 'plain/text'}
    }).then((res) => {
      return res.json()
    }).then(res => {
      console.log('RESPONSE', res)
      setManifest(res)
      if (!haveManifest(res)) setManifests([res, ...manifests])
      setLoading(false)
    })
  }

  const haveManifest = (man) => {
    return _.some(manifests, m => m.id === man.id)
  }

  return (
    <div className="Panel">
      <Current manifest={manifest} loading={loading} makeRequest={getManifest}/>
      <History manifests={manifests} showManifest={showManifest}/>
    </div>
  );
}
