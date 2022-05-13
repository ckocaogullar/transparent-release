// Copyright 2022 The Project Oak Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wrappers

import (
	"io/ioutil"
	"testing"
)

const testRekorLogPath = "experimental/auth-logic/test_data/rekor_entry.json"
const testPubKeyPath = "experimental/auth-logic/test_data/oak_ec_public.pem"
const testUnexpiredEndorsementFilePath = "experimental/auth-logic/test_data/oak_endorsement.json"
const rekorPublicKeyPath = "experimental/auth-logic/test_data/rekor_public_key.pem"

func TestRekorLogWrapper(t *testing.T) {
	rekorLogEntryBytes, err := ioutil.ReadFile(testRekorLogPath)
	if err != nil {
		t.Errorf("could not read rekor log file %v\n", testRekorLogPath)
	}

	prodTeamKeyBytes, err := ioutil.ReadFile(testPubKeyPath)
	if err != nil {
		t.Errorf("could not parse prod team pub key from file: %s", testPubKeyPath)
	}

	endorsementBytes, err := ioutil.ReadFile(testUnexpiredEndorsementFilePath)
	if err != nil {
		t.Errorf("could not read endorsement file: %s, %v", testUnexpiredEndorsementFilePath, err)
	}

	rekorKeyBytes, err := ioutil.ReadFile(rekorPublicKeyPath)
	if err != nil {
		t.Errorf("could not parse rekord pub key from file: %s", rekorKeyBytes)
	}

	// Test of VerifyRekordEntry
	err = VerifyRekorEntry(rekorLogEntryBytes, prodTeamKeyBytes, rekorKeyBytes, endorsementBytes)
	if err != nil {
		t.Errorf("rekord entry verification should have succeeded for this test: %v", err)
	}

	// ---- Test of RekorWrapper
	// Expected output of wrapper:
	want := `RekorLogCheck says {
hasValidBodySignature("oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::RekorLogEntry").
hasValidInclusionProof("oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::RekorLogEntry").
signerIsProductTeam("oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::RekorLogEntry").
contentsMatch("oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::RekorLogEntry", "oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::EndorsementFile").
"oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::EndorsementFile" canActAs ValidRekorEntry :- hasValidBodySignature("oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::RekorLogEntry"), hasValidInclusionProof("oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::RekorLogEntry"), hasCorrectPubKey("oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::RekorLogEntry"), contentsMatch("oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::RekorLogEntry", "oak_functions_loader:0f2189703c57845e09d8ab89164a4041c0af0a62::EndorsementFile").
}`

	testRekorLogWrapper := RekorLogWrapper{
		rekorLogEntryBytes:  rekorLogEntryBytes,
		productTeamKeyBytes: prodTeamKeyBytes,
		rekorPublicKeyBytes: rekorKeyBytes,
		endorsementBytes:    endorsementBytes,
	}

	rekorLogStatement, err := EmitStatementAs(Principal{Contents: "RekorLogCheck"}, testRekorLogWrapper)
	if err != nil {
		t.Errorf("couldn't get rekor log statement: %v, %v", testRekorLogWrapper, err)
	}

	got := rekorLogStatement.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s\n", got, want)
	}

}

func TestVerifySignedEntryTimestamp(t *testing.T) {
	rekorLogEntryBytes, err := ioutil.ReadFile(testRekorLogPath)
	if err != nil {
		t.Errorf("could not read rekor log file %v\n", testRekorLogPath)
	}

	rekorKeyBytes, err := ioutil.ReadFile(rekorPublicKeyPath)
	if err != nil {
		t.Errorf("could not parse rekord pub key from file: %s", rekorKeyBytes)
	}

	logEntryAnon, err := getLogEntryAnonFromBytes(rekorLogEntryBytes)
	if err != nil {
		t.Errorf("could not get logEntryAnon: %v", err)
	}

	err = verifySignedEntryTimestamp(logEntryAnon, rekorKeyBytes)
	if err != nil {
		t.Fatalf("could not verify signed entry timestamp: %v", err)
	}

}
