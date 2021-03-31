ethereum.autoRefreshOnNetworkChange = false

// setup Web3 provider
const network = 'rinkeby'
const provider = new ethers.getDefaultProvider(network, {etherscan: '5TC41U8EV1GMIQKE3VV7CX1FJT1NE1BF6V'})

const init = () => {
    document.getElementById('auth').addEventListener('click', async () => {
      await auth()
    }, true)
    document.getElementById('refresh').addEventListener('click', async () => {
      await refresh()
    }, true)
}

const auth = async () => {
  // get challenge from API while using the current account from Metamask
  let req = new Request(`http://localhost:3000/login/${web3.eth.defaultAccount}`, {
    method: 'GET'
  })
  let res = await fetch(req)
  let data = await res.json()
  let sig = ""

  // make sure we're signing from the right account
  if (web3.eth.defaultAccount === data.address) {
    sig = await sign(data.challenge)
  }

  // reply with signed response
  req = new Request(`http://localhost:3000/login/${web3.eth.defaultAccount}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: `{"signature": "${sig}"}`
  })
  res = await fetch(req)
  data = await res.json();
  if (res.status === 200) {
    /// all good, think about storing the authz token
    allow(data.token)
  } else if (res.status === 403) {
    // forbidden
    deny(data)
  } else {
    throw new Error('Something went wrong on api server!');
  }
}

// refresh token
const refresh = async () => {
  req = new Request(`http://localhost:3000/refresh`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${document.getElementById('token').innerText}`
    }
  })
  res = await fetch(req)
  data = await res.json();
  if (res.status === 200) {
    allow(data.token)
  } else if (res.status === 403) {
    deny(data)
  } else {
    throw new Error('Something went wrong on api server!');
  }
}

// sign message with Metamask
const sign = (msg) => {
  const signer = (new ethers.providers.Web3Provider(window.ethereum)).getSigner()
  return signer.signMessage(msg)
}

// update DOM when allowing access
const allow = (token) => {
  document.getElementById('not_allowed').hidden = true
  document.getElementById('allowed').hidden = false
  document.getElementById('token').innerText = token
}

// update DOM when denying access
const deny = (name) => {
  // name = (!name) ? '<em>Unregistered ENS name</em>' : name
  document.getElementById('allowed').hidden = true
  const p = document.createElement('p')
  p.style.color = 'red'
  p.innerHTML = `Not allowed`
  document.getElementById('not_allowed').appendChild(p)
  document.getElementById('not_allowed').hidden = false
}

init()
