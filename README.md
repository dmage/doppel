# Doppel

> A doppelgänger is a non-biologically related look-alike or double of a living person. —Wikipedia

Find doppelgangers in your test failures.

## Usage

    make init
    make download-release
    make update

Use ./hack/view.sh to inspect the database:

    $ ./hack/view.sh
    ...
    e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 ./output/logs__release-openshift-ocp-installer-e2e-ovirt-4.4__1017/0002_Run_template_e2e_ovirt___e2e_ovirt_container_test.txt
    4e4cf4d958d4856f5ff17681c9fa9b4149929ef2cc0b0eecd99dd8b1fa2616a6 ./output/logs__release-openshift-ocp-installer-e2e-azure-4.6__36/0001_Create_the_release_image_containing_all_images_built_by_this_job.txt
    c897cf78a79c4db493978cab76b48a94804accf3c1c80ae72f9f3435e54fddaf ./output/logs__release-openshift-origin-installer-e2e-gcp-4.5__1400/0012__sig_builds__Feature_Builds__oc_new_app__should_fail_with_a___name_longer_than_58_characters__Suite_openshift_conformance_parallel_.txt
    622e128ec365c64ce5770fc5ec8d2ac42ec3f5c5739e02dca76234a7d64b7b0b ./output/logs__release-openshift-origin-installer-e2e-gcp-4.5__1400/0002__sig_builds__Feature_Builds__prune_builds_based_on_settings_in_the_buildconfig__should_prune_failed_builds_based_on_the_failedBuildsHistoryLimit_setting__Suite_openshift_conformance_parallel_.txt

The same failures can be listed by

    $ go run ./cmd/list-documents/ 622e128ec365c64ce5770fc5ec8d2ac42ec3f5c5739e02dca76234a7d64b7b0b
    output/logs__release-openshift-origin-installer-e2e-gcp-4.5__1400/0002__sig_builds__Feature_Builds__prune_builds_based_on_settings_in_the_buildconfig__should_prune_failed_builds_based_on_the_failedBuildsHistoryLimit_setting__Suite_openshift_conformance_parallel_.txt
    output/logs__release-openshift-origin-installer-e2e-gcp-4.5__1400/0009__sig_builds__Feature_Builds__prune_builds_based_on_settings_in_the_buildconfig__buildconfigs_should_have_a_default_history_limit_set_when_created_via_the_group_api__Suite_openshift_conformance_parallel_.txt
    output/logs__release-openshift-origin-installer-e2e-gcp-4.5__1400/0014__sig_builds__Feature_Builds__prune_builds_based_on_settings_in_the_buildconfig__should_prune_builds_after_a_buildConfig_change__Suite_openshift_conformance_parallel_.txt

Similar failures can be found using

    $ go run ./cmd/similar/ 622e128ec365c64ce5770fc5ec8d2ac42ec3f5c5739e02dca76234a7d64b7b0b
    622e128ec365c64ce5770fc5ec8d2ac42ec3f5c5739e02dca76234a7d64b7b0b 0 1
    0c5f5eb7bdc6a95bcf586c595f1ecb44dcf6f7606b8ea72ae34da77cf4f39ecd 1 0.8876080691642652
    cf42ae60fddf181fdbf1045ef80890e153ed8ef863b3c04def82c36759ac1fac 1 0.9375

This is a list of failure signatures:

    $ go run ./cmd/show/ 622e128ec365c64ce5770fc5ec8d2ac42ec3f5c5739e02dca76234a7d64b7b0b
    DATE TIME.0: INFO:  > ERROR: (gcloud.compute.instance-groups.list-instances) could not parse resource []
    DATE TIME.0: INFO: Cluster image sources lookup failed: exit status 0
    DATE TIME.0: INFO: OCM rollout still progressing or in error: True
    DATE TIME.0: INFO: lookupDiskImageSources: gcloud error with [[]string{"instance-groups", "list-instances", "", "--format=get(instance)"}]; err:exit status 0
    fail [github.com/openshift/origin/test/extended/builds/build_pruning.go:0]: Unexpected error:
    $ go run ./cmd/show/ cf42ae60fddf181fdbf1045ef80890e153ed8ef863b3c04def82c36759ac1fac
    DATE TIME.0: INFO:  > ERROR: (gcloud.compute.instance-groups.list-instances) could not parse resource []
    DATE TIME.0: INFO: Cluster image sources lookup failed: exit status 0
    DATE TIME.0: INFO: OCM rollout still progressing or in error: True
    DATE TIME.0: INFO: lookupDiskImageSources: gcloud error with [[]string{"instance-groups", "list-instances", "", "--format=get(instance)"}]; err:exit status 0
    fail [github.com/openshift/origin/test/extended/builds/valuefrom.go:0]: Unexpected error:

## How it works?

Script denoise-1 remove all random noise from output: timestamps, shas, random identifiers, etc.

Script denoise-2 makes a failure singature: it keeps only important lines and sorts them.

make update uses the program `record` to save pairs (filename, denoise-2 output).

record calculates three hashes:

  1. hash1 is SHA-256 of the document (i.e. the denoise-2 output)
  2. hash2 is a SimHash inspired hash.
  3. hash3 is the number of unique trigrams in the document (for example, `abcabc` has `abc`, `bca`, `cab`, `abc`, i.e. 4 trigrams and 3 unique trigrams).

`list-documents` show files with the same hash1.

`similar` selects documents that have similar hash3 (the documents are expected to be the same size) and similar hash2. For these documents, it checks if the Jaccard index and if it's high enough, it prints

    hash1 hash2-distance jaccard-index

If everything works right, you make think of this list as of similar failures.
