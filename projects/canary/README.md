CANARY - HACKATHON SUBMISSION

# Canary

- **Track(s):** Public Goods, Advanced Cryptography (MPC/Threshold Crypto), Censorship Resistance, 
- **Team/Contributors:** Kieran P, Lina P, Ryan C
- **Repository:** https://github.com/TheThirdRoom/canary
- **Demo:** https://canary-demo-z4ftr.ondigitalocean.app/

## Description (TL;DR)
Canary is an automated deadman switch for truth protection, using TACo threshold encryption, decentralized storage, and smart contracts to ensure sensitive information reaches the public even if journalists, activists, or whistleblowers are silenced. Users create encrypted "dossiers" that automatically release to trusted recipients (or the public) if they fail to check in, providing uncensorable protection for those who protect democracy.

## Problem
Journalists, activists, and whistleblowers face life-threatening risks when exposing corruption and injustice, with 1,700+ journalists killed since 2006 (UN stat). Current solutions can be compromised through coercion, server shutdowns, or organizational pressure. Critical stories and evidence are lost forever when truth-tellers are silenced, leaving corruption unexposed and justice unserved.

## Solution
Canary automates truth protection through decentralized deadman switches that survive any attempt at suppression. Users upload encrypted files to Codex or IPFS, set custom check-in schedules, and designate trusted recipients (or the public at large). If check-ins stop, TACo's threshold encryption automatically releases decryption keys to recipients, ensuring stories publish regardless of what happens to the original user.

## Technology Stack
Frontend: React with TypeScript, Tailwind CSS, Web3 wallet integration
Storage: Codex (could use any decentralized storage layer)
Encryption: TACo threshold encryption for conditional access control, client-side local file encryption, e2ee

## Privacy Impact
Canary uses TACo's threshold encryption to ensure no single point of failure, while Codex provides censorship-resistant storage. Users maintain complete control over their data with client-side encryption.

## Real-World Use Cases

Journalist in a war zone who wants to auto-release their encrypted story if detained or worse...

Activist who wants to release private keys of crypto funds for bail. (possible inheritance angle?)

Crossing a border and want to share your story only if something goes wrong.

A crypto hodler who is traveling to France and afraid of getting kidnapped so you want close friends to know your geo-location (coming soon™️)


## Business Logic
Enterprise (Newsrooms, NGO's Lawfirms) and Direct Consumer (Crypto holders, Activists, OSINT community, whistleblowers)


## What's Next
Keep building of course, adding features such as geo location, live audio and photo documentation, obfuscating the crypto part of the app, collecting feedback from potential end users.
