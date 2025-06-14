# ACTivate Groups

- **Track(s):** Censorship Resistance, Applied Encryption
- **Team/Contributors:** Bálint @bosi95, József @Kexort, András @aranyia; from @Solar-Punk-Ltd
- **Repository:** https://github.com/Solar-Punk-Ltd/hackathon-2025-berlin-submissions
- **Demo:** [Link to live demo, video, or screenshots]

## Description (TL;DR)

ACTivate Groups is a decentralized, censorship-resistant information sharing application for creating and managing
private, self-organizing groups. Built on Swarm, it uses public key cryptography to ensure that only authorized members
can participate, providing a secure communication channel shielded from surveillance and disruption.

## Problem

Activists, journalists, and other at-risk groups often rely on centralized communication platforms that are vulnerable
to government surveillance, censorship, and single points of failure. These platforms can be blocked, shut down, or
forced to reveal user data, putting sensitive communications and the safety of individuals at risk. There is a critical
need for a communication tool that is not dependent on centralized infrastructure and prioritizes user privacy,
anonymity, and operational resilience.

## Solution

ACTivate Groups addresses this by leveraging the decentralized storage and peer-to-peer architecture of Swarm. Our
solution has two key layers of protection:

1. **Group Access Control:** We use Swarm's Access Control Trie (ACT) as a decentralized and tamper-proof membership
   list. A group administrator controls access by adding or removing members' Swarm public keys. Only wallets with a
   corresponding private key can access the group's content on the network. For the initial version, this allows for a
   one-way communication model where the admin shares information with the group. Two-way communication can be achieved
   by users creating their own reciprocal groups and cross-adding each other as members.

2. **End-to-End Encryption:** All messages are end-to-end encrypted when uploaded to Swarm with ACT. This provides a
   crucial defense-in-depth layer. The encryption mechanism of ACT then ensures that this content is further protected
   and only decipherable by a group member with the specific decryption key designated for them.

This layered approach creates a robust, secure, and pseudonymous environment for groups to communicate and coordinate
safely.

## Technology Stack

- **Core Infrastructure:** [Swarm](https://www.ethswarm.org/) (for decentralized storage and p2p communication)

- **Swarm Interaction:** Bee
  Client ([swarm-mobile](https://github.com/Solar-Punk-Ltd/swarm-mobile)), [Gnosis chain](https://www.gnosischain.com/)

- **Group Management:** Swarm's Access Control
  Trie ([ACT](https://solarpunk.buzz/introducing-the-access-control-trie-act-in-swarm/))

- **Mobile App:** [Go Mobile](https://pkg.go.dev/golang.org/x/mobile) for building a cross-platform mobile app,
  [Fyne](https://fyne.io/) for UI development

## Privacy Impact

This project significantly enhances user privacy in several ways:

- **Pseudonymity:** Users are identified only by their Swarm public key, with no link to real-world identities.

- **No Central Servers:** All data is fragmented and distributed across the Swarm network. There are no central servers
  to attack, subpoena, or surveil.

- **Metadata Minimization:** The decentralized architecture resists traffic analysis and the mapping of social graphs,
  as messages are routed through the Swarm network rather than point-to-point.

- **Censorship Resistance:** By design, the distributed network is resilient to takedowns or blocks by central
  authorities.

## Real-World Use Cases

- **Activists & Organizers:** Groups operating in high-risk environments can coordinate actions without fear of their
  communication channels being shut down or monitored.

- **Journalists & Whistleblowers:** Provides a secure channel to communicate with sources, protecting their anonymity
  and the integrity of the information.

- **Private Communities:** Any group seeking a truly private and sovereign space for communication, free from the
  data-mining practices of mainstream platforms.

## Business Logic

As a public good, the primary goal of ACTivate Groups is not monetization but sustainability. The main operational cost
is funding the Swarm "postage stamps" required to ensure data persistence. Our sustainability model could include:

- **Group Treasuries:** Implementing a feature for groups to pool funds (e.g., in xBZZ) to collectively purchase and
  top-up postage stamps for their chat history.

- **Grant Funding:** Seeking support from foundations dedicated to digital freedom, privacy, and open-source technology.

- **Donations:** Accepting donations from users and supporters who believe in the project's mission.

## What's Next

- **Multi-Admin & Decentralized Governance:** Evolving from a single-admin model to allow for multi-signature control
  over the group's ACT, or even more complex decentralized governance mechanisms.

- **Automated Postage Stamp Management:** Abstracting away the complexity of managing postage stamps by creating
  automated, group-funded top-ups to improve UX and prevent data loss.

- **Enhanced Deniability:** Implementing features like scheduled message deletion and improving resistance to metadata
  analysis to provide users with plausible deniability.

- **Improved Onboarding:** Creating a seamless out-of-band process for securely exchanging the initial public keys
  needed to join a group.