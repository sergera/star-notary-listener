# Star Notary Listener

<p>this service listens to events from the <a href="https://github.com/sergera/star-notary">Star Notary</a> smart contract<br>
and communicates with the star-notary-backend (to be announced)</p>

## Install

<pre><code>make install</pre></code>

## Run

<pre><code>make run</pre></code>

## Requirements

<p>have <a href="https://go.dev/">Go</a> installed and binary added to PATH</p>

<p>register at <a href="https://infura.io/">Infura</a> or another RPC provider that exposes websocket endpoints</p>

## Environment Variables

<p>locally, environment variables are declared in 'application.conf'</p>

###### RPC_PROVIDER_WEBSOCKET_URL:

<p>websocket url to deployed network provided by Infura or chosen RPC provider</p>

###### CONTRACT_ADDRESS:

<p>address of currently deployed smart contract</p>

###### CONFIRMATION_BLOCKS:

<p>number (integer) of confirmation blocks before an event is considered cannon</p>

###### CONFIRMATION_SLEEP_SECONDS:

<p>number (integer) of seconds that the service waits between RPC provider calls for event confirmation</p>

###### LOG_PATH (optional):

<p>full path to log directory</p>
<p>if not provided logs to project root directory</p>

## Go Contract Creation

<pre><code>make contract</pre></code>

<p>to create a go contract, the following variables must be set in the Makefile:</p>

###### TRUFFLE_PROJECT_ROOT_PATH:

<p>full path to root directory of truffle project containing the solidity (.sol) smart contract file</p>

###### SOLIDITY_VERSION:

<p>version of solidity used in the smart contract</p>

###### CONTRACT_NAME:

<p>name of the solidity (.sol) smart contract file</p>
