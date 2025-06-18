# ProofLeak

- **Track(s):** [1 - Advanced Cryptography and 5 - Public Goods]
- **Team/Contributors:** [steelfeel.eth - zk, blockchain, solidity; 0xxdana.eth - UI/UX, PM, GTM]
- **Repository:** [ProofLeak Repository](https://github.com/steel-feel/hackathon-2025-berlin-submissions.git)
- **Demo:** [ProofLeak Presentation](https://pitch.com/v/proofleak---a-secure-platform-to-whistleblow-gxgkp7)

## Description (TL;DR)
A private and secure whistleblowing platform using zk email verification and onchain attestations to ensure source credibility without compromising anonymity.

We want to build a cryptographic infrastructure to free the truth. With ZK Proofs we don't have to choose between anonimity and credibility anymore. We can have both. We can have truth.

## Problem
Most anonymous reporting platforms rely on central servers, lack cryptographic verification, and expose whistleblowers to surveillance or metadata leaks. Without decentralized trust or ZK-based guarantees, both source safety and information integrity are at risk.

## Solution
Using zero-knowledge email proofs, we verify source legitimacy (e.g., domain or role-based credentials) without exposing identity. Blockchain attestations anchor evidence onchain with tamper-proof audit trails. This enables transparency: sources can selectively disclose credentials while maintaining full anonymity.

## Technology Stack
noir ,IPFS, zkemail, TypeScript, JavaScript, Rust, React, zkverify

## Privacy Impact
This replaces institutional trust with cryptographic proof. ZKPs allow email domain and DKIM verification without revealing message content or identity. No metadata leaks, no storage of sensitive information. Whistleblowers remain fully private, while communities gain confidence in the credibility of what is submitted.

## Real-World Use Cases
- Corporate insiders can prove access to internal systems or domains without revealing identity
- Verified legal and journalist-only forums can receive sensitive documents anonymously
- Newsrooms can cryptographically validate anonymous sources
- Vulnerable individuals can safely disclose important information without fear of retaliation

## Business Logic
We believe this tool should be accessible to all, a Common Public Good funded by aligned communities.

## What's Next
- Technical community validation
- Develop MVP for full submission flow (wallet + eml + onchain attestation)
- Expand group attestation features for verified journalist and legal forums -> the possibility to create closed forums where individuals meeting the criteria can receive anonymous, verfied leaks in order to take the necessary steps until the wider public can be made aware of. Sometimes information needs to be private until it can be made public for the greater good. 
- Launch with transparency-focused partners and fund as a public good
