package spec_types

import (
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/merge"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/beacon/sharding"
	"github.com/protolambda/ztyp/view"
)

type SpecType struct {
	Alloc   func(spec *common.Spec) common.SSZObj
	TypeDef func(spec *common.Spec) view.TypeDef
}

var Phase0SpecTypes = map[string]SpecType{
	"BeaconState":             {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.BeaconState)) }, func(spec *common.Spec) view.TypeDef { return phase0.BeaconStateType(spec) }},
	"BeaconBlock":             {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.BeaconBlock)) }, func(spec *common.Spec) view.TypeDef { return phase0.BeaconBlockType(spec) }},
	"BeaconBlockHeader":       {func(spec *common.Spec) common.SSZObj { return new(common.BeaconBlockHeader) }, func(spec *common.Spec) view.TypeDef { return common.BeaconBlockHeaderType }},
	"SignedBeaconBlock":       {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.SignedBeaconBlock)) }, func(spec *common.Spec) view.TypeDef { return phase0.BeaconStateType(spec) }},
	"SignedBeaconBlockHeader": {func(spec *common.Spec) common.SSZObj { return new(common.SignedBeaconBlockHeader) }, func(spec *common.Spec) view.TypeDef { return common.SignedBeaconBlockHeaderType }},
	"BeaconBlockBody":         {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.BeaconBlockBody)) }, func(spec *common.Spec) view.TypeDef { return phase0.BeaconBlockBodyType(spec) }},

	"AttestationData":     {func(spec *common.Spec) common.SSZObj { return new(phase0.AttestationData) }, func(spec *common.Spec) view.TypeDef { return phase0.AttestationDataType }},
	"Attestation":         {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.Attestation)) }, func(spec *common.Spec) view.TypeDef { return phase0.AttestationType(spec) }},
	"AttesterSlashing":    {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.AttesterSlashing)) }, func(spec *common.Spec) view.TypeDef { return phase0.AttesterSlashingType(spec) }},
	"ProposerSlashing":    {func(spec *common.Spec) common.SSZObj { return new(phase0.ProposerSlashing) }, func(spec *common.Spec) view.TypeDef { return phase0.ProposerSlashingType }},
	"Deposit":             {func(spec *common.Spec) common.SSZObj { return new(common.Deposit) }, func(spec *common.Spec) view.TypeDef { return common.DepositType }},
	"DepositData":         {func(spec *common.Spec) common.SSZObj { return new(common.DepositData) }, func(spec *common.Spec) view.TypeDef { return common.DepositDataType }},
	"DepositMessage":      {func(spec *common.Spec) common.SSZObj { return new(common.DepositMessage) }, func(spec *common.Spec) view.TypeDef { return common.DepositMessageType }},
	"VoluntaryExit":       {func(spec *common.Spec) common.SSZObj { return new(phase0.VoluntaryExit) }, func(spec *common.Spec) view.TypeDef { return phase0.VoluntaryExitType }},
	"SignedVoluntaryExit": {func(spec *common.Spec) common.SSZObj { return new(phase0.SignedVoluntaryExit) }, func(spec *common.Spec) view.TypeDef { return phase0.SignedVoluntaryExitType }},
	"Eth1Data":            {func(spec *common.Spec) common.SSZObj { return new(common.Eth1Data) }, func(spec *common.Spec) view.TypeDef { return common.Eth1DataType }},
	"ForkData":            {func(spec *common.Spec) common.SSZObj { return new(common.ForkData) }, func(spec *common.Spec) view.TypeDef { return common.ForkDataType }},
	"Fork":                {func(spec *common.Spec) common.SSZObj { return new(common.Fork) }, func(spec *common.Spec) view.TypeDef { return common.ForkType }},
	"Checkpoint":          {func(spec *common.Spec) common.SSZObj { return new(common.Checkpoint) }, func(spec *common.Spec) view.TypeDef { return common.CheckpointType }},
	"Validator":           {func(spec *common.Spec) common.SSZObj { return new(phase0.Validator) }, func(spec *common.Spec) view.TypeDef { return phase0.ValidatorType }},
	"IndexedAttestation":  {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.IndexedAttestation)) }, func(spec *common.Spec) view.TypeDef { return phase0.IndexedAttestationType(spec) }},
	"PendingAttestation":  {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.PendingAttestation)) }, func(spec *common.Spec) view.TypeDef { return phase0.PendingAttestationType(spec) }},

	"SigningData":    {func(spec *common.Spec) common.SSZObj { return new(common.SigningData) }, func(spec *common.Spec) view.TypeDef { return common.SigningDataType }},
	"Slot":           {func(spec *common.Spec) common.SSZObj { return new(common.Slot) }, func(spec *common.Spec) view.TypeDef { return common.SlotType }},
	"Epoch":          {func(spec *common.Spec) common.SSZObj { return new(common.Epoch) }, func(spec *common.Spec) view.TypeDef { return common.EpochType }},
	"CommitteeIndex": {func(spec *common.Spec) common.SSZObj { return new(common.CommitteeIndex) }, func(spec *common.Spec) view.TypeDef { return common.CommitteeIndexType }},
	"ValidatorIndex": {func(spec *common.Spec) common.SSZObj { return new(common.ValidatorIndex) }, func(spec *common.Spec) view.TypeDef { return common.ValidatorIndexType }},
	"Gwei":           {func(spec *common.Spec) common.SSZObj { return new(common.Gwei) }, func(spec *common.Spec) view.TypeDef { return common.GweiType }},
	"Root":           {func(spec *common.Spec) common.SSZObj { return new(common.Root) }, func(spec *common.Spec) view.TypeDef { return view.RootType }},
	"Hash32":         {func(spec *common.Spec) common.SSZObj { return new(common.Hash32) }, func(spec *common.Spec) view.TypeDef { return common.Hash32Type }},
	"Version":        {func(spec *common.Spec) common.SSZObj { return new(common.Version) }, func(spec *common.Spec) view.TypeDef { return common.VersionType }},
	"DomainType":     {func(spec *common.Spec) common.SSZObj { return new(common.BLSDomainType) }, func(spec *common.Spec) view.TypeDef { return common.BLSDomainTypeTreeType }},
	"ForkDigest":     {func(spec *common.Spec) common.SSZObj { return new(common.ForkDigest) }, func(spec *common.Spec) view.TypeDef { return common.ForkDigestType }},
	"Domain":         {func(spec *common.Spec) common.SSZObj { return new(common.BLSDomain) }, func(spec *common.Spec) view.TypeDef { return common.BLSDomainTreeType }},
	"BLSPubkey":      {func(spec *common.Spec) common.SSZObj { return new(common.BLSPubkey) }, func(spec *common.Spec) view.TypeDef { return common.BLSPubkeyType }},
	"BLSSignature":   {func(spec *common.Spec) common.SSZObj { return new(common.BLSSignature) }, func(spec *common.Spec) view.TypeDef { return common.BLSSignatureType }},
}

var AltairSpecTypes = map[string]SpecType{
	"BeaconState":             {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.BeaconState)) }, func(spec *common.Spec) view.TypeDef { return altair.BeaconStateType(spec) }},
	"BeaconBlock":             {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.BeaconBlock)) }, func(spec *common.Spec) view.TypeDef { return altair.BeaconBlockType(spec) }},
	"BeaconBlockHeader":       {func(spec *common.Spec) common.SSZObj { return new(common.BeaconBlockHeader) }, func(spec *common.Spec) view.TypeDef { return common.BeaconBlockHeaderType }},
	"SignedBeaconBlock":       {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SignedBeaconBlock)) }, func(spec *common.Spec) view.TypeDef { return altair.SignedBeaconBlockType(spec) }},
	"SignedBeaconBlockHeader": {func(spec *common.Spec) common.SSZObj { return new(common.SignedBeaconBlockHeader) }, func(spec *common.Spec) view.TypeDef { return common.SignedBeaconBlockHeaderType }},
	"BeaconBlockBody":         {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.BeaconBlockBody)) }, func(spec *common.Spec) view.TypeDef { return altair.BeaconBlockBodyType(spec) }},

	"AttestationData":     {func(spec *common.Spec) common.SSZObj { return new(phase0.AttestationData) }, func(spec *common.Spec) view.TypeDef { return phase0.AttestationDataType }},
	"Attestation":         {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.Attestation)) }, func(spec *common.Spec) view.TypeDef { return phase0.AttestationType(spec) }},
	"AttesterSlashing":    {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.AttesterSlashing)) }, func(spec *common.Spec) view.TypeDef { return phase0.AttesterSlashingType(spec) }},
	"ProposerSlashing":    {func(spec *common.Spec) common.SSZObj { return new(phase0.ProposerSlashing) }, func(spec *common.Spec) view.TypeDef { return phase0.ProposerSlashingType }},
	"Deposit":             {func(spec *common.Spec) common.SSZObj { return new(common.Deposit) }, func(spec *common.Spec) view.TypeDef { return common.DepositType }},
	"DepositData":         {func(spec *common.Spec) common.SSZObj { return new(common.DepositData) }, func(spec *common.Spec) view.TypeDef { return common.DepositDataType }},
	"DepositMessage":      {func(spec *common.Spec) common.SSZObj { return new(common.DepositMessage) }, func(spec *common.Spec) view.TypeDef { return common.DepositMessageType }},
	"VoluntaryExit":       {func(spec *common.Spec) common.SSZObj { return new(phase0.VoluntaryExit) }, func(spec *common.Spec) view.TypeDef { return phase0.VoluntaryExitType }},
	"SignedVoluntaryExit": {func(spec *common.Spec) common.SSZObj { return new(phase0.SignedVoluntaryExit) }, func(spec *common.Spec) view.TypeDef { return phase0.SignedVoluntaryExitType }},
	"Eth1Data":            {func(spec *common.Spec) common.SSZObj { return new(common.Eth1Data) }, func(spec *common.Spec) view.TypeDef { return common.Eth1DataType }},
	"ForkData":            {func(spec *common.Spec) common.SSZObj { return new(common.ForkData) }, func(spec *common.Spec) view.TypeDef { return common.ForkDataType }},
	"Fork":                {func(spec *common.Spec) common.SSZObj { return new(common.Fork) }, func(spec *common.Spec) view.TypeDef { return common.ForkType }},
	"Checkpoint":          {func(spec *common.Spec) common.SSZObj { return new(common.Checkpoint) }, func(spec *common.Spec) view.TypeDef { return common.CheckpointType }},
	"Validator":           {func(spec *common.Spec) common.SSZObj { return new(phase0.Validator) }, func(spec *common.Spec) view.TypeDef { return phase0.ValidatorType }},
	"IndexedAttestation":  {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.IndexedAttestation)) }, func(spec *common.Spec) view.TypeDef { return phase0.IndexedAttestationType(spec) }},

	"SigningData":    {func(spec *common.Spec) common.SSZObj { return new(common.SigningData) }, func(spec *common.Spec) view.TypeDef { return common.SigningDataType }},
	"Slot":           {func(spec *common.Spec) common.SSZObj { return new(common.Slot) }, func(spec *common.Spec) view.TypeDef { return common.SlotType }},
	"Epoch":          {func(spec *common.Spec) common.SSZObj { return new(common.Epoch) }, func(spec *common.Spec) view.TypeDef { return common.EpochType }},
	"CommitteeIndex": {func(spec *common.Spec) common.SSZObj { return new(common.CommitteeIndex) }, func(spec *common.Spec) view.TypeDef { return common.CommitteeIndexType }},
	"ValidatorIndex": {func(spec *common.Spec) common.SSZObj { return new(common.ValidatorIndex) }, func(spec *common.Spec) view.TypeDef { return common.ValidatorIndexType }},
	"Gwei":           {func(spec *common.Spec) common.SSZObj { return new(common.Gwei) }, func(spec *common.Spec) view.TypeDef { return common.GweiType }},
	"Root":           {func(spec *common.Spec) common.SSZObj { return new(common.Root) }, func(spec *common.Spec) view.TypeDef { return view.RootType }},
	"Hash32":         {func(spec *common.Spec) common.SSZObj { return new(common.Hash32) }, func(spec *common.Spec) view.TypeDef { return common.Hash32Type }},
	"Version":        {func(spec *common.Spec) common.SSZObj { return new(common.Version) }, func(spec *common.Spec) view.TypeDef { return common.VersionType }},
	"DomainType":     {func(spec *common.Spec) common.SSZObj { return new(common.BLSDomainType) }, func(spec *common.Spec) view.TypeDef { return common.BLSDomainTypeTreeType }},
	"ForkDigest":     {func(spec *common.Spec) common.SSZObj { return new(common.ForkDigest) }, func(spec *common.Spec) view.TypeDef { return common.ForkDigestType }},
	"Domain":         {func(spec *common.Spec) common.SSZObj { return new(common.BLSDomain) }, func(spec *common.Spec) view.TypeDef { return common.BLSDomainTreeType }},
	"BLSPubkey":      {func(spec *common.Spec) common.SSZObj { return new(common.BLSPubkey) }, func(spec *common.Spec) view.TypeDef { return common.BLSPubkeyType }},
	"BLSSignature":   {func(spec *common.Spec) common.SSZObj { return new(common.BLSSignature) }, func(spec *common.Spec) view.TypeDef { return common.BLSSignatureType }},

	"LightClientSnapshot": {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.LightClientSnapshot)) }, func(spec *common.Spec) view.TypeDef { return altair.LightClientSnapshotType(spec) }},
	"LightClientUpdate":   {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.LightClientUpdate)) }, func(spec *common.Spec) view.TypeDef { return altair.LightClientUpdateType(spec) }},

	"SyncAggregate":               {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SyncAggregate)) }, func(spec *common.Spec) view.TypeDef { return altair.SyncAggregateType(spec) }},
	"SyncAggregatorSelectionData": {func(spec *common.Spec) common.SSZObj { return new(altair.SyncAggregatorSelectionData) }, func(spec *common.Spec) view.TypeDef { return altair.SyncAggregatorSelectionDataType }},
	"SyncCommitteeContribution":   {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SyncCommitteeContribution)) }, func(spec *common.Spec) view.TypeDef { return altair.SyncCommitteeContributionType(spec) }},
	"ContributionAndProof":        {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.ContributionAndProof)) }, func(spec *common.Spec) view.TypeDef { return altair.ContributionAndProofType(spec) }},
	"SignedContributionAndProof":  {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SignedContributionAndProof)) }, func(spec *common.Spec) view.TypeDef { return altair.SignedContributionAndProofType(spec) }},
	"SyncCommitteeMessage":        {func(spec *common.Spec) common.SSZObj { return new(altair.SyncCommitteeMessage) }, func(spec *common.Spec) view.TypeDef { return altair.SyncCommitteeMessageType }},
	"SyncCommittee":               {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(common.SyncCommittee)) }, func(spec *common.Spec) view.TypeDef { return common.SyncCommitteeType(spec) }},
}

var MergeSpecTypes = map[string]SpecType{
	"BeaconState":             {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(merge.BeaconState)) }, func(spec *common.Spec) view.TypeDef { return merge.BeaconStateType(spec) }},
	"BeaconBlock":             {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(merge.BeaconBlock)) }, func(spec *common.Spec) view.TypeDef { return merge.BeaconBlockType(spec) }},
	"BeaconBlockHeader":       {func(spec *common.Spec) common.SSZObj { return new(common.BeaconBlockHeader) }, func(spec *common.Spec) view.TypeDef { return common.BeaconBlockHeaderType }},
	"SignedBeaconBlock":       {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(merge.SignedBeaconBlock)) }, func(spec *common.Spec) view.TypeDef { return merge.SignedBeaconBlockType(spec) }},
	"SignedBeaconBlockHeader": {func(spec *common.Spec) common.SSZObj { return new(common.SignedBeaconBlockHeader) }, func(spec *common.Spec) view.TypeDef { return common.SignedBeaconBlockHeaderType }},
	"BeaconBlockBody":         {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(merge.BeaconBlockBody)) }, func(spec *common.Spec) view.TypeDef { return merge.BeaconBlockBodyType(spec) }},

	"AttestationData":     {func(spec *common.Spec) common.SSZObj { return new(phase0.AttestationData) }, func(spec *common.Spec) view.TypeDef { return phase0.AttestationDataType }},
	"Attestation":         {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.Attestation)) }, func(spec *common.Spec) view.TypeDef { return phase0.AttestationType(spec) }},
	"AttesterSlashing":    {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.AttesterSlashing)) }, func(spec *common.Spec) view.TypeDef { return phase0.AttesterSlashingType(spec) }},
	"ProposerSlashing":    {func(spec *common.Spec) common.SSZObj { return new(phase0.ProposerSlashing) }, func(spec *common.Spec) view.TypeDef { return phase0.ProposerSlashingType }},
	"Deposit":             {func(spec *common.Spec) common.SSZObj { return new(common.Deposit) }, func(spec *common.Spec) view.TypeDef { return common.DepositType }},
	"DepositData":         {func(spec *common.Spec) common.SSZObj { return new(common.DepositData) }, func(spec *common.Spec) view.TypeDef { return common.DepositDataType }},
	"DepositMessage":      {func(spec *common.Spec) common.SSZObj { return new(common.DepositMessage) }, func(spec *common.Spec) view.TypeDef { return common.DepositMessageType }},
	"VoluntaryExit":       {func(spec *common.Spec) common.SSZObj { return new(phase0.VoluntaryExit) }, func(spec *common.Spec) view.TypeDef { return phase0.VoluntaryExitType }},
	"SignedVoluntaryExit": {func(spec *common.Spec) common.SSZObj { return new(phase0.SignedVoluntaryExit) }, func(spec *common.Spec) view.TypeDef { return phase0.SignedVoluntaryExitType }},
	"Eth1Data":            {func(spec *common.Spec) common.SSZObj { return new(common.Eth1Data) }, func(spec *common.Spec) view.TypeDef { return common.Eth1DataType }},
	"ForkData":            {func(spec *common.Spec) common.SSZObj { return new(common.ForkData) }, func(spec *common.Spec) view.TypeDef { return common.ForkDataType }},
	"Fork":                {func(spec *common.Spec) common.SSZObj { return new(common.Fork) }, func(spec *common.Spec) view.TypeDef { return common.ForkType }},
	"Checkpoint":          {func(spec *common.Spec) common.SSZObj { return new(common.Checkpoint) }, func(spec *common.Spec) view.TypeDef { return common.CheckpointType }},
	"Validator":           {func(spec *common.Spec) common.SSZObj { return new(phase0.Validator) }, func(spec *common.Spec) view.TypeDef { return phase0.ValidatorType }},
	"IndexedAttestation":  {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.IndexedAttestation)) }, func(spec *common.Spec) view.TypeDef { return phase0.IndexedAttestationType(spec) }},

	"SigningData":    {func(spec *common.Spec) common.SSZObj { return new(common.SigningData) }, func(spec *common.Spec) view.TypeDef { return common.SigningDataType }},
	"Slot":           {func(spec *common.Spec) common.SSZObj { return new(common.Slot) }, func(spec *common.Spec) view.TypeDef { return common.SlotType }},
	"Epoch":          {func(spec *common.Spec) common.SSZObj { return new(common.Epoch) }, func(spec *common.Spec) view.TypeDef { return common.EpochType }},
	"CommitteeIndex": {func(spec *common.Spec) common.SSZObj { return new(common.CommitteeIndex) }, func(spec *common.Spec) view.TypeDef { return common.CommitteeIndexType }},
	"ValidatorIndex": {func(spec *common.Spec) common.SSZObj { return new(common.ValidatorIndex) }, func(spec *common.Spec) view.TypeDef { return common.ValidatorIndexType }},
	"Gwei":           {func(spec *common.Spec) common.SSZObj { return new(common.Gwei) }, func(spec *common.Spec) view.TypeDef { return common.GweiType }},
	"Root":           {func(spec *common.Spec) common.SSZObj { return new(common.Root) }, func(spec *common.Spec) view.TypeDef { return view.RootType }},
	"Hash32":         {func(spec *common.Spec) common.SSZObj { return new(common.Hash32) }, func(spec *common.Spec) view.TypeDef { return common.Hash32Type }},
	"Version":        {func(spec *common.Spec) common.SSZObj { return new(common.Version) }, func(spec *common.Spec) view.TypeDef { return common.VersionType }},
	"DomainType":     {func(spec *common.Spec) common.SSZObj { return new(common.BLSDomainType) }, func(spec *common.Spec) view.TypeDef { return common.BLSDomainTypeTreeType }},
	"ForkDigest":     {func(spec *common.Spec) common.SSZObj { return new(common.ForkDigest) }, func(spec *common.Spec) view.TypeDef { return common.ForkDigestType }},
	"Domain":         {func(spec *common.Spec) common.SSZObj { return new(common.BLSDomain) }, func(spec *common.Spec) view.TypeDef { return common.BLSDomainTreeType }},
	"BLSPubkey":      {func(spec *common.Spec) common.SSZObj { return new(common.BLSPubkey) }, func(spec *common.Spec) view.TypeDef { return common.BLSPubkeyType }},
	"BLSSignature":   {func(spec *common.Spec) common.SSZObj { return new(common.BLSSignature) }, func(spec *common.Spec) view.TypeDef { return common.BLSSignatureType }},

	"LightClientSnapshot": {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.LightClientSnapshot)) }, func(spec *common.Spec) view.TypeDef { return altair.LightClientSnapshotType(spec) }},
	"LightClientUpdate":   {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.LightClientUpdate)) }, func(spec *common.Spec) view.TypeDef { return altair.LightClientUpdateType(spec) }},

	"SyncAggregate":               {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SyncAggregate)) }, func(spec *common.Spec) view.TypeDef { return altair.SyncAggregateType(spec) }},
	"SyncAggregatorSelectionData": {func(spec *common.Spec) common.SSZObj { return new(altair.SyncAggregatorSelectionData) }, func(spec *common.Spec) view.TypeDef { return altair.SyncAggregatorSelectionDataType }},
	"SyncCommitteeContribution":   {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SyncCommitteeContribution)) }, func(spec *common.Spec) view.TypeDef { return altair.SyncCommitteeContributionType(spec) }},
	"ContributionAndProof":        {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.ContributionAndProof)) }, func(spec *common.Spec) view.TypeDef { return altair.ContributionAndProofType(spec) }},
	"SignedContributionAndProof":  {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SignedContributionAndProof)) }, func(spec *common.Spec) view.TypeDef { return altair.SignedContributionAndProofType(spec) }},
	"SyncCommitteeMessage":        {func(spec *common.Spec) common.SSZObj { return new(altair.SyncCommitteeMessage) }, func(spec *common.Spec) view.TypeDef { return altair.SyncCommitteeMessageType }},
	"SyncCommittee":               {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(common.SyncCommittee)) }, func(spec *common.Spec) view.TypeDef { return common.SyncCommitteeType(spec) }},

	"ExecutionPayload":       {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(common.ExecutionPayload)) }, func(spec *common.Spec) view.TypeDef { return common.ExecutionPayloadType(spec) }},
	"ExecutionPayloadHeader": {func(spec *common.Spec) common.SSZObj { return new(common.ExecutionPayloadHeader) }, func(spec *common.Spec) view.TypeDef { return common.ExecutionPayloadHeaderType }},
	"PayloadTransactions":    {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(common.PayloadTransactions)) }, func(spec *common.Spec) view.TypeDef { return common.PayloadTransactionsType(spec) }},
	"LogsBloom":              {func(spec *common.Spec) common.SSZObj { return new(common.LogsBloom) }, func(spec *common.Spec) view.TypeDef { return common.LogsBloomType }},
}

var ShardingSpecTypes = map[string]SpecType{
	"BeaconState":             {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.BeaconState)) }, func(spec *common.Spec) view.TypeDef { return sharding.BeaconStateType(spec) }},
	"BeaconBlock":             {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.BeaconBlock)) }, func(spec *common.Spec) view.TypeDef { return sharding.BeaconBlockType(spec) }},
	"BeaconBlockHeader":       {func(spec *common.Spec) common.SSZObj { return new(common.BeaconBlockHeader) }, func(spec *common.Spec) view.TypeDef { return common.BeaconBlockHeaderType }},
	"SignedBeaconBlock":       {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.SignedBeaconBlock)) }, func(spec *common.Spec) view.TypeDef { return sharding.SignedBeaconBlockType(spec) }},
	"SignedBeaconBlockHeader": {func(spec *common.Spec) common.SSZObj { return new(common.SignedBeaconBlockHeader) }, func(spec *common.Spec) view.TypeDef { return common.SignedBeaconBlockHeaderType }},
	"BeaconBlockBody":         {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.BeaconBlockBody)) }, func(spec *common.Spec) view.TypeDef { return sharding.BeaconBlockBodyType(spec) }},

	"AttestationData":  {func(spec *common.Spec) common.SSZObj { return new(sharding.AttestationData) }, func(spec *common.Spec) view.TypeDef { return sharding.AttestationDataType }},
	"Attestation":      {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.Attestation)) }, func(spec *common.Spec) view.TypeDef { return sharding.AttestationType(spec) }},
	"AttesterSlashing": {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.AttesterSlashing)) }, func(spec *common.Spec) view.TypeDef { return sharding.AttesterSlashingType(spec) }},

	"ProposerSlashing":    {func(spec *common.Spec) common.SSZObj { return new(phase0.ProposerSlashing) }, func(spec *common.Spec) view.TypeDef { return phase0.ProposerSlashingType }},
	"Deposit":             {func(spec *common.Spec) common.SSZObj { return new(common.Deposit) }, func(spec *common.Spec) view.TypeDef { return common.DepositType }},
	"DepositData":         {func(spec *common.Spec) common.SSZObj { return new(common.DepositData) }, func(spec *common.Spec) view.TypeDef { return common.DepositDataType }},
	"DepositMessage":      {func(spec *common.Spec) common.SSZObj { return new(common.DepositMessage) }, func(spec *common.Spec) view.TypeDef { return common.DepositMessageType }},
	"VoluntaryExit":       {func(spec *common.Spec) common.SSZObj { return new(phase0.VoluntaryExit) }, func(spec *common.Spec) view.TypeDef { return phase0.VoluntaryExitType }},
	"SignedVoluntaryExit": {func(spec *common.Spec) common.SSZObj { return new(phase0.SignedVoluntaryExit) }, func(spec *common.Spec) view.TypeDef { return phase0.SignedVoluntaryExitType }},
	"Eth1Data":            {func(spec *common.Spec) common.SSZObj { return new(common.Eth1Data) }, func(spec *common.Spec) view.TypeDef { return common.Eth1DataType }},
	"ForkData":            {func(spec *common.Spec) common.SSZObj { return new(common.ForkData) }, func(spec *common.Spec) view.TypeDef { return common.ForkDataType }},
	"Fork":                {func(spec *common.Spec) common.SSZObj { return new(common.Fork) }, func(spec *common.Spec) view.TypeDef { return common.ForkType }},
	"Checkpoint":          {func(spec *common.Spec) common.SSZObj { return new(common.Checkpoint) }, func(spec *common.Spec) view.TypeDef { return common.CheckpointType }},
	"Validator":           {func(spec *common.Spec) common.SSZObj { return new(phase0.Validator) }, func(spec *common.Spec) view.TypeDef { return phase0.ValidatorType }},
	"IndexedAttestation":  {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(phase0.IndexedAttestation)) }, func(spec *common.Spec) view.TypeDef { return phase0.IndexedAttestationType(spec) }},

	"SigningData":    {func(spec *common.Spec) common.SSZObj { return new(common.SigningData) }, func(spec *common.Spec) view.TypeDef { return common.SigningDataType }},
	"Slot":           {func(spec *common.Spec) common.SSZObj { return new(common.Slot) }, func(spec *common.Spec) view.TypeDef { return common.SlotType }},
	"Epoch":          {func(spec *common.Spec) common.SSZObj { return new(common.Epoch) }, func(spec *common.Spec) view.TypeDef { return common.EpochType }},
	"CommitteeIndex": {func(spec *common.Spec) common.SSZObj { return new(common.CommitteeIndex) }, func(spec *common.Spec) view.TypeDef { return common.CommitteeIndexType }},
	"ValidatorIndex": {func(spec *common.Spec) common.SSZObj { return new(common.ValidatorIndex) }, func(spec *common.Spec) view.TypeDef { return common.ValidatorIndexType }},
	"Gwei":           {func(spec *common.Spec) common.SSZObj { return new(common.Gwei) }, func(spec *common.Spec) view.TypeDef { return common.GweiType }},
	"Root":           {func(spec *common.Spec) common.SSZObj { return new(common.Root) }, func(spec *common.Spec) view.TypeDef { return view.RootType }},
	"Hash32":         {func(spec *common.Spec) common.SSZObj { return new(common.Hash32) }, func(spec *common.Spec) view.TypeDef { return common.Hash32Type }},
	"Version":        {func(spec *common.Spec) common.SSZObj { return new(common.Version) }, func(spec *common.Spec) view.TypeDef { return common.VersionType }},
	"DomainType":     {func(spec *common.Spec) common.SSZObj { return new(common.BLSDomainType) }, func(spec *common.Spec) view.TypeDef { return common.BLSDomainTypeTreeType }},
	"ForkDigest":     {func(spec *common.Spec) common.SSZObj { return new(common.ForkDigest) }, func(spec *common.Spec) view.TypeDef { return common.ForkDigestType }},
	"Domain":         {func(spec *common.Spec) common.SSZObj { return new(common.BLSDomain) }, func(spec *common.Spec) view.TypeDef { return common.BLSDomainTreeType }},
	"BLSPubkey":      {func(spec *common.Spec) common.SSZObj { return new(common.BLSPubkey) }, func(spec *common.Spec) view.TypeDef { return common.BLSPubkeyType }},
	"BLSSignature":   {func(spec *common.Spec) common.SSZObj { return new(common.BLSSignature) }, func(spec *common.Spec) view.TypeDef { return common.BLSSignatureType }},

	"LightClientSnapshot": {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.LightClientSnapshot)) }, func(spec *common.Spec) view.TypeDef { return altair.LightClientSnapshotType(spec) }},
	"LightClientUpdate":   {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.LightClientUpdate)) }, func(spec *common.Spec) view.TypeDef { return altair.LightClientUpdateType(spec) }},

	"SyncAggregate":               {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SyncAggregate)) }, func(spec *common.Spec) view.TypeDef { return altair.SyncAggregateType(spec) }},
	"SyncAggregatorSelectionData": {func(spec *common.Spec) common.SSZObj { return new(altair.SyncAggregatorSelectionData) }, func(spec *common.Spec) view.TypeDef { return altair.SyncAggregatorSelectionDataType }},
	"SyncCommitteeContribution":   {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SyncCommitteeContribution)) }, func(spec *common.Spec) view.TypeDef { return altair.SyncCommitteeContributionType(spec) }},
	"ContributionAndProof":        {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.ContributionAndProof)) }, func(spec *common.Spec) view.TypeDef { return altair.ContributionAndProofType(spec) }},
	"SignedContributionAndProof":  {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(altair.SignedContributionAndProof)) }, func(spec *common.Spec) view.TypeDef { return altair.SignedContributionAndProofType(spec) }},
	"SyncCommitteeMessage":        {func(spec *common.Spec) common.SSZObj { return new(altair.SyncCommitteeMessage) }, func(spec *common.Spec) view.TypeDef { return altair.SyncCommitteeMessageType }},
	"SyncCommittee":               {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(common.SyncCommittee)) }, func(spec *common.Spec) view.TypeDef { return common.SyncCommitteeType(spec) }},

	"ExecutionPayload":       {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(common.ExecutionPayload)) }, func(spec *common.Spec) view.TypeDef { return common.ExecutionPayloadType(spec) }},
	"ExecutionPayloadHeader": {func(spec *common.Spec) common.SSZObj { return new(common.ExecutionPayloadHeader) }, func(spec *common.Spec) view.TypeDef { return common.ExecutionPayloadHeaderType }},
	"PayloadTransactions":    {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(common.PayloadTransactions)) }, func(spec *common.Spec) view.TypeDef { return common.PayloadTransactionsType(spec) }},
	"LogsBloom":              {func(spec *common.Spec) common.SSZObj { return new(common.LogsBloom) }, func(spec *common.Spec) view.TypeDef { return common.LogsBloomType }},

	"Shard":        {func(spec *common.Spec) common.SSZObj { return new(common.Shard) }, func(spec *common.Spec) view.TypeDef { return common.ShardType }},
	"BuilderIndex": {func(spec *common.Spec) common.SSZObj { return new(common.BuilderIndex) }, func(spec *common.Spec) view.TypeDef { return common.BuilderIndexType }},
	"Builder":      {func(spec *common.Spec) common.SSZObj { return new(sharding.Builder) }, func(spec *common.Spec) view.TypeDef { return sharding.BuilderType }},

	"BLSPoint":                 {func(spec *common.Spec) common.SSZObj { return new(common.BLSPoint) }, func(spec *common.Spec) view.TypeDef { return common.BLSPointType }},
	"DataCommitment":           {func(spec *common.Spec) common.SSZObj { return new(sharding.DataCommitment) }, func(spec *common.Spec) view.TypeDef { return sharding.DataCommitmentType }},
	"AttestedDataCommitment":   {func(spec *common.Spec) common.SSZObj { return new(sharding.AttestedDataCommitment) }, func(spec *common.Spec) view.TypeDef { return sharding.AttestedDataCommitmentType }},
	"ShardData":                {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.ShardData)) }, func(spec *common.Spec) view.TypeDef { return sharding.ShardDataType(spec) }},
	"ShardBlobBody":            {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.ShardBlobBody)) }, func(spec *common.Spec) view.TypeDef { return sharding.ShardBlobBodyType(spec) }},
	"ShardBlobBodySummary":     {func(spec *common.Spec) common.SSZObj { return new(sharding.ShardBlobBodySummary) }, func(spec *common.Spec) view.TypeDef { return sharding.ShardBlobBodySummaryType }},
	"ShardBlobHeader":          {func(spec *common.Spec) common.SSZObj { return new(sharding.ShardBlobHeader) }, func(spec *common.Spec) view.TypeDef { return sharding.ShardBlobHeaderType }},
	"SignedShardBlobHeader":    {func(spec *common.Spec) common.SSZObj { return new(sharding.SignedShardBlobHeader) }, func(spec *common.Spec) view.TypeDef { return sharding.SignedShardBlobHeaderType }},
	"ShardBlob":                {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.ShardBlob)) }, func(spec *common.Spec) view.TypeDef { return sharding.ShardBlobType(spec) }},
	"SignedShardBlob":          {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.SignedShardBlob)) }, func(spec *common.Spec) view.TypeDef { return sharding.SignedShardBlobType(spec) }},
	"PendingShardHeader":       {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.PendingShardHeader)) }, func(spec *common.Spec) view.TypeDef { return sharding.PendingShardHeaderType(spec) }},
	"ShardBlobReference":       {func(spec *common.Spec) common.SSZObj { return new(sharding.ShardBlobReference) }, func(spec *common.Spec) view.TypeDef { return sharding.ShardBlobReferenceType }},
	"SignedShardBlobReference": {func(spec *common.Spec) common.SSZObj { return new(sharding.SignedShardBlobReference) }, func(spec *common.Spec) view.TypeDef { return sharding.SignedShardBlobReferenceType }},
	"ShardProposerSlashing":    {func(spec *common.Spec) common.SSZObj { return new(sharding.ShardProposerSlashing) }, func(spec *common.Spec) view.TypeDef { return sharding.ShardProposerSlashingType }},
	"ShardWork":                {func(spec *common.Spec) common.SSZObj { return spec.Wrap(new(sharding.ShardWork)) }, func(spec *common.Spec) view.TypeDef { return sharding.ShardWorkType(spec) }},
}

var TypesByPhase = map[string]map[string]SpecType{
	"phase0":   Phase0SpecTypes,
	"altair":   AltairSpecTypes,
	"merge":    MergeSpecTypes,
	"sharding": ShardingSpecTypes,
}

var Phases = []string{"phase0", "altair", "merge", "sharding"}

func TypeNames(types map[string]SpecType) []string {
	out := make([]string, 0, len(types))
	for k := range types {
		out = append(out, k)
	}
	return out
}
