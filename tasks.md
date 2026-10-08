# Jira Tickets

Generated from `jira_links` on `2026-05-12`.

## Status

| Key | Title | Status |
|-----|-------|--------|
| SDK-1929 | [Terraform] Change 'WaitUntilAvailable' behavior for DBaaS cluster creation | DONE |
| SDK-2528 | [Terraform] Modify test.yml workflow | DONE |

> **Status column instructions:** When Claude works on a ticket, it fills this column with `DONE` if the task was completed, or `BLOCKED` if it could not be finished (append a brief reason in parentheses, e.g. `BLOCKED (missing API access)`). Leave blank until worked on.

---

## SDK-1929 — [Terraform] Change 'WaitUntilAvailable' behavior for DBaaS cluster creation

**Link:** https://hosting-jira.1and1.org/browse/SDK-1929

**Description:**

Change ***WaitUntilAvailable*** behavior for DBaaS cluster creation.

When creating a DBaaS cluster, the cluster can enter the ***FAILED*** state. After entering this state, there is no point in waiting for the cluster to become ***AVAILABLE*** since this will never happen, so we need to change the ***Wait*** mechanism.

**Attachments:**

*(none)*

---

## SDK-2528 — [Terraform] Modify test.yml workflow

**Link:** https://hosting-jira.1and1.org/browse/SDK-2528

**Description:**

From **test.yml**:

```yaml
 name: Run tests for ${{ github.event.inputs.test-tags }}
  if: ${{ github.event.inputs.test-tags != '' && github.events.inputs.failfast == true }}
  run: go test ./... -v -failfast -timeout 240m -tags ${{ github.event.inputs.test-tags }}
- name: Run tests for ${{ github.event.inputs.test-tags }} no failfast
  if: ${{ github.event.inputs.test-tags != '' && github.events.inputs.failfast == false }}
  run: go test ./... -v -timeout 6h -tags ${{ github.event.inputs.test-tags }}
- name: Run tests without tags
  if: ${{ github.event.inputs.test-tags == ''}}
  run: go test ./... -v -failfast -timeout 180m
```

Timeout value differs when using the **failfast** option. When setting the **failfast** option, I don't expect a change in the timeout value. The timeout value should be the same for both scenarios but it would be nice to be able to specify the timeout value as an input parameter.

**Attachments:**

*(none)*
