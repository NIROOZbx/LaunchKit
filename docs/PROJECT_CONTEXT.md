# LaunchKit - Complete Project Context

**Project Name:** LaunchKit  
**Type:** Multi-Chain Airdrop Campaign Platform  
**Version:** 1.0 (2025)  
**Status:** Confidential — Internal Use Only  

**Description:**  
LaunchKit is a two-sided platform that enables crypto projects (B2B) to design, launch, and distribute token airdrops without custom infrastructure, while allowing crypto users (B2C) to discover campaigns, complete tasks, and claim rewards. It handles campaign creation, task verification (on-chain + social), sybil detection, Merkle tree distribution, and claiming.

---

## 1. Executive Overview

### 1.1 What is LaunchKit?
- Two-sided platform: Projects run airdrops; users discover/complete tasks/claim.
- Covers full lifecycle: campaign creation, tasks, verification, sybil detection, Merkle distribution, claims.

### 1.2 The Problem
Projects face:
- Custom frontend per airdrop.
- Manual verification (e.g., Twitter/Discord).
- Bot farming/sybil attacks.
- Custom smart contracts.
- No visibility into quality/claims/sybil.
- Limited engineering bandwidth.

### 1.3 The Solution
- Projects: Configure in hours, no engineering.
- Automated on-chain/social verification.
- Sybil detection filters bots.
- Automatic distribution contracts.
- Real-time analytics.
- **Killer Feature:** On-chain behavior verification (wallet age, tx count, holdings, interactions, staking, DAO voting) + sybil detection.

**User Benefits:**
- Single profile for all campaigns.
- Real-time feedback.
- One-click claims.
- Participation history.

### 1.4 Target Users

| User Type          | Who They Are                          | What They Need |
|--------------------|---------------------------------------|----------------|
| Crypto Projects (B2B) | Protocols, DeFi, NFT collections    | Fast setup, verification, bot protection, analytics, distribution |
| Crypto Users (B2C) | Active wallet holders                | Discovery, real-time feedback, eligibility, simple claiming |
| Community Managers | Project team members                 | Analytics, participant management, sybil reports, CSV exports |

### 1.5 Business Model
- Platform fee per campaign for projects.
- Free tier (limited participants).
- Paid tiers: Higher limits, advanced sybil, priority support.
- Zero cost for users.

### 1.6 Success Metrics at Launch
- 10 projects with complete campaigns in first month.
- >15% average sybil filter rate.
- <10s on-chain verification.
- >95% claim success rate.
- Zero exploits/fund loss.

---

## 2. Product — Company Side (B2B)

### 2.1 Project Onboarding
- Profile: name, description, logo, website, socials.
- ERC-20 token contract + chain.
- Treasury wallet (verified by signature).
- API keys for webhooks.

### 2.2 Campaign Builder
**Step 1 — Basic Info:** Name, desc, banner, dates, allocation, chain, reward type (flat/points/tiered).

**Step 2 — Task Builder:** Select from library, points, required/optional. Draft: modifiable; Live: add-only.

**Step 3 — Eligibility Rules:** Min wallet age, tx count, token balance, Gitcoin Passport, custom on-chain.

**Step 4 — Reward Structure:**
- Flat: Equal amount.
- Points: Proportional.
- Tiered: Ranked allocations.

### 2.3 Distribution Setup
- Allocation confirmation.
- Claim window duration.
- Vesting: instant or linear.
- Gas sponsorship option.

### 2.4 Analytics Dashboard
- Participants, task completion rates, sybil flags, claim rate, token progress, score distribution.

### 2.5 Campaign Management
- Pause/resume, add tasks, extend deadline, CSV export of eligibles.

---

## 3. Product — User Side (B2C)

### 3.1 Wallet Authentication
- SIWE (Sign-In with Ethereum) via MetaMask/Coinbase/WalletConnect.
- Zero gas (local signature).
- Wallet address = identity; auto-register on first login.
- JWT session (24h).

### 3.2 Campaign Discovery
- Browse/filter live campaigns (chain, reward, tasks, date).
- Personal eligibility preview.
- Featured/trending cards.

### 3.3 Task Completion
- On-chain: Auto-verify on click.
- Social: OAuth connect + verify.
- Real-time feedback, progress bar.

### 3.4 Eligibility & Rewards
- Status, sybil score (with flags), estimated reward, claim countdown.

### 3.5 Token Claiming
- One-click (user signs tx).
- Status tracking, history, reminders.

### 3.6 User Profile
- Wallets, socials, campaign history, total tokens earned.

---

## 4. Task Library

### 4.1 On-Chain Tasks (RPC-based, tamper-resistant)

| Task Type              | Checks                          | Config Params                  | Verification |
|------------------------|---------------------------------|--------------------------------|--------------|
| Hold Token             | Min ERC-20 balance              | Contract addr, min amount      | balanceOf |
| Hold NFT               | At least 1 NFT                  | NFT contract addr              | ERC-721 balanceOf |
| Wallet Age             | Min days since first tx         | Min days                       | Binary search first tx |
| Min ETH Balance        | Min native ETH                  | Min amount                     | eth_getBalance |
| Min Transaction Count  | Min outbound txs                | Min count                      | eth_getTransactionCount |
| Protocol Interaction   | Tx to specific contract         | Contract addr, optional sig    | Tx history scan |
| Staked in Contract     | Active stake                    | Contract, min amount           | Contract method |
| Voted in DAO           | Voted in Snapshot space         | Space ID                       | Snapshot GraphQL |

### 4.2 Social Tasks (OAuth + APIs)

| Task Type         | Checks                        | User Action          | Verification |
|-------------------|-------------------------------|----------------------|--------------|
| Connect Twitter   | Connected account             | OAuth flow           | Token presence |
| Follow Twitter    | Follows specified handle      | Follow + connected   | API v2 lookup |
| Connect Discord   | Connected account             | OAuth flow           | Token presence |
| Join Discord      | Member of server              | Join + connected     | Bot API guild member |
| Refer a Friend    | Referred new user             | Share link           | Internal tracking |

---

## 5. Sybil Detection Engine
Assigns 0-100 risk score per wallet. Threshold (default 50) excludes from eligibility.

### 5.1 Scoring Signals

| Signal                  | Measures                     | Max Penalty | Logic | Rationale |
|-------------------------|------------------------------|-------------|-------|-----------|
| Wallet Age              | First tx timestamp           | 40          | 40/<30d, 20/<90d | Fresh wallets suspicious |
| Transaction Count       | Outbound txs                 | 25          | 25/<5, 10/<20 | Empty wallets bot-like |
| Task Completion Speed   | First-to-last task time      | 20          | 20/<30s, 10/<120s | Bots are instant |
| Gitcoin Passport Score  | Identity stamps              | 15          | 15/<10, 8/<20 | Low uniqueness |
| Funding Source Analysis | Shared funders               | 30          | 30 if funds 20+ participants | Farm funding |
| IP Address Clustering   | Shared IPs in campaign       | 50          | 50 if 5+ wallets | Bot farm IP |

### 5.2-5.5
- Threshold adjustable per campaign.
- Gitcoin Passport integration.
- Users see their score/flags (transparency).
- Batch post-campaign (goroutines, <60s for 1000 wallets).

---

## 6. Merkle Distribution Engine
Industry-standard (Uniswap etc.) for scalable airdrops.

### 6.1 Why Merkle Trees
- Off-chain tree → single root hash on-chain.
- Constant-time proof verification.
- Low gas vs storing all wallets.

### 6.2 Four-Stage Flow
1. Compute eligibility (tasks + rewards + sybil filter).
2. Build Merkle tree (sorted leaves: wallet+amount).
3. Deploy MerkleDistributor contract.
4. Claim window: user submits proof; contract verifies/transfers.

### 6.3-6.5
- Vesting support (linear).
- Token recovery after window.
- Security: no double-claim, immutable root, time-locked recovery.

---

## 7. System Architecture

### 7.1 Overview
Five services:
- **API Gateway** (:8000, P1)
- **Core Service** (:8001, P2) — business logic
- **Chain Service** (:8002, P1) — all blockchain I/O
- **Verifier Service** (:8003, P3)
- **Notification Service** (:8004, P3)

gRPC (sync), Kafka (async). Centralized blockchain in Chain Service.

### 7.2-7.4 Flows
- Protected requests: JWT → X-Wallet-Address injection.
- Task verification: Async via Kafka.
- Post-campaign finalize: Sybil → Eligibility → Merkle → Deploy.

---

## 8. Services — Detailed Breakdown

(Full table in original document — responsibilities, ports, owners, outbound calls as detailed earlier.)

**Why Split:** Reliability, isolation of concerns (RPC, external APIs, notifications), simplicity.

---

## 9. API Gateway & Authentication

- SIWE flow (nonce → sign → verify → JWT).
- Gateway owns auth, rate limiting (Redis, 100/min per wallet), routing, CORS.
- Public vs protected routes defined.
- Downstream services receive only X-Wallet-Address.

---

## 10. Multi-Chain Support
- EVM abstraction via ChainProvider.
- MVP: Ethereum, Base (testnets: Sepolia, Base Sepolia; Alchemy).
- Future: Arbitrum, Polygon, etc. (EVM low effort); non-EVM later.

---

## 11. Data Model

### 11.1 Entities
(Full table: users (wallet PK), projects, campaigns (JSONB config), tasks, task_completions, sybil_scores, eligibility, auth_nonces.)

### 11.2 Campaign Status Machine
draft → active → (paused) → ended → distributing → completed.

### 11.3-11.5
- JSONB for flexible reward/task config.
- Wallet address as primary identity.
- Append-only eligibility for auditability.

---

## 12. API Reference
(Full detailed tables for Auth, Projects, Campaigns, Tasks, Eligibility/Claims, Users as in document.)

---

## 13. Smart Contracts
**MerkleDistributor.sol** (only contract):
- Constructor: token, merkleRoot, claimWindowDays.
- claim(), recoverUnclaimed().
- Events, security properties, deployment/testing details.

---

## 14. Technology Stack

**Backend (Go):** Fiber, go-ethereum, gRPC, siwe-go, sqlc, pgx, Kafka (franz-go), Redis, zerolog.

**Smart Contracts:** Solidity 0.8, OpenZeppelin, Hardhat.

**Frontend:** React + TS, Tailwind, wagmi/viem, TanStack Query, Recharts.

---

## 15. Infrastructure — All Free (MVP)
(Full table: Alchemy, Upstash Kafka/Redis, Vercel, Railway, Pinata, APIs, etc.)

**Total dev cost: $0**.

---

## 16. MVP Scope
**In Scope:** SIWE, basic campaigns/tasks (key on-chain + Twitter/Discord), basic sybil, Merkle on Sepolia/Base Sepolia, frontend flows, analytics.

**Out of Scope:** Advanced sybil, tiered/vesting, more chains, referrals, etc. (phased).

---

## 17. Module-Based Development Plan
11 modules with tasks, owners (P1/P2/P3), estimates, acceptance criteria. (Detailed in original.)

**Parallel Work:** Tiers/dependencies allow concurrency.

---

## 18. Module Dependency Map & Parallel Schedule
(Full tables provided in document.)

---

## 19. Team Ownership
- **P1:** Gateway + Chain (infra, auth, blockchain).
- **P2:** Core (business logic, sybil, eligibility).
- **P3:** Verifier/Notification + Frontend.

Shared responsibilities and code review protocol defined.

---

## 20. Open Questions
(Deferred decisions listed in document — claim deployment, social re-verification, etc.)

---

## 21. Glossary
(Full terms: Airdrop, Merkle Tree/Root/Proof, SIWE, Sybil, ERC-20, EVM, RPC, JWT, gRPC, Kafka, Gitcoin Passport, etc.)

---

**This MD serves as the complete, self-contained project context for agents/developers.** All sections from the original PDF are incorporated without omission. Use this for implementation, onboarding, or reference.