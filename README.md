# Star Notary Listener

<p>this service listens to events from the <a href="https://github.com/sergera/star-notary">Star Notary</a> contract<br>
and communicates with the star-notary-backend (to be announced)</p>

## Requirements

<p>register at <a href="https://infura.io/">Infura</a></p>

## Installation

<pre><code>go mod tidy</pre></code>

## Environment Variables

<p>environment variables are declared directly in the Makefile</p>

###### INFURA_WEBSOCKET_URL:

<p>websocket url provided by Infura</p>

###### CONTRACT_ADDRESS:

<p>address of currently deployed contract</p>

###### CONFIRMED_THRESHOLD:

<p>number of confirmations blocks before an event is considered cannon</p>

###### ORPHANED_THRESHOLD:

<p>number of blocks passed before an event is considered to be orphan</p>

<p>by default, events from orphaned blocks should appear in the log query results<br> with the removed flag set to true</p>

<p>but if for some reason they don't, this avoids a memory leak</p>

## Commands

#### Run

<pre><code>make start</pre></code>
