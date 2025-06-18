# MagicENS

- Track(s): User Onboarding - Applied Encryption  
- Team/Contributors: zkfriendly, sembrestels, FaezehShakouri  
- Repository: https://github.com/stealtheth  
- Demo: https://www.youtube.com/watch?v=Hv_N-B20C2s

## Description (TL;DR)

MagicENS adds a privacy layer to ENS names.  
It allows the owner to share their name once with the world, and anyone — using any wallet — can send them money privately without installing new software or visiting any new web pages beyond their current wallet interface.

Behind the scenes, the ENS name resolves to a different address each time it is used in a wallet. Essentially, it is like generating a brand new receiving address every time — automatically.

## Problem

Currently, ENS names do not provide any privacy for the receiver.  
If you share your ENS name with someone, they can see how much money you’ve received and from whom. This creates a privacy issue and leaks sensitive information — such as who you transact with and how much you’ve received.

## Solution

We use stealth addresses to generate a new receiving address each time — all controlled by your main account, which remains unrevealed.  
On the ENS side, we use an off-chain resolver that handles generating and resolving these addresses.

## Technology Stack

- Semaphore  
- Viem  
- Fluidkey  
- Stealth addresses  
- ERC-4337 paymaster  
- ZeroDev accounts

## Privacy Impact

Hides the activity and interactions of an ENS name on the public blockchain.

## Real-World Use Cases

Anyone with an ENS name can use this to add privacy to their address — without requiring the sender to install anything or change their normal wallet behavior.

## Business Logic

[Sustainability/monetization considerations]

## What's Next

- Deposit into privacy pools  
- Finish Semaphore paymaster integrations  
- Work toward ENS standardizations to remove the need for a server  
- Implement an easy-to-use withdrawal mechanism  
- Add aggregated stealth addresses for improved usability
