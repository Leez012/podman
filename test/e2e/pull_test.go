package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/containers/podman/v2/test/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Podman pull", func() {
	var (
		tempdir    string
		err        error
		podmanTest *PodmanTestIntegration
	)

	BeforeEach(func() {
		tempdir, err = CreateTempDirInTempDir()
		if err != nil {
			os.Exit(1)
		}
		podmanTest = PodmanTestCreate(tempdir)
		podmanTest.Setup()
	})

	AfterEach(func() {
		podmanTest.Cleanup()
		f := CurrentGinkgoTestDescription()
		processTestResult(f)

	})

	It("podman pull from docker a not existing image", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "ibetthisdoesntexistthere:foo"})
		session.WaitWithDefaultTimeout()
		Expect(session).To(ExitWithError())
	})

	It("podman pull from docker with tag", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "busybox:glibc"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.PodmanNoCache([]string{"rmi", "busybox:glibc"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull from docker without tag", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "busybox"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.PodmanNoCache([]string{"rmi", "busybox"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull from alternate registry with tag", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", nginx})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.PodmanNoCache([]string{"rmi", nginx})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull from alternate registry without tag", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "quay.io/libpod/alpine_nginx"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.PodmanNoCache([]string{"rmi", "quay.io/libpod/alpine_nginx"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull by digest", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "alpine@sha256:1072e499f3f655a032e88542330cf75b02e7bdf673278f701d7ba61629ee3ebe"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.PodmanNoCache([]string{"rmi", "alpine:none"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull by digest (image list)", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "--override-arch=arm64", ALPINELISTDIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		// inspect using the digest of the list
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINELISTDIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(HavePrefix("[]"))
		// inspect using the digest of the list
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINELISTDIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTDIGEST))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINEARM64DIGEST))
		// inspect using the digest of the arch-specific image's manifest
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINEARM64DIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(HavePrefix("[]"))
		// inspect using the digest of the arch-specific image's manifest
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINEARM64DIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTDIGEST))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINEARM64DIGEST))
		// inspect using the image ID
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINEARM64ID})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(HavePrefix("[]"))
		// inspect using the image ID
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINEARM64ID})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTDIGEST))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINEARM64DIGEST))
		// remove using the digest of the list
		session = podmanTest.PodmanNoCache([]string{"rmi", ALPINELISTDIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull by instance digest (image list)", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "--override-arch=arm64", ALPINEARM64DIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		// inspect using the digest of the list
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINELISTDIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Not(Equal(0)))
		// inspect using the digest of the list
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINELISTDIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Not(Equal(0)))
		// inspect using the digest of the arch-specific image's manifest
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINEARM64DIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(HavePrefix("[]"))
		// inspect using the digest of the arch-specific image's manifest
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINEARM64DIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(Not(ContainSubstring(ALPINELISTDIGEST)))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINEARM64DIGEST))
		// inspect using the image ID
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINEARM64ID})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(HavePrefix("[]"))
		// inspect using the image ID
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINEARM64ID})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(Not(ContainSubstring(ALPINELISTDIGEST)))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINEARM64DIGEST))
		// remove using the digest of the instance
		session = podmanTest.PodmanNoCache([]string{"rmi", ALPINEARM64DIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull by tag (image list)", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "--override-arch=arm64", ALPINELISTTAG})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		// inspect using the tag we used for pulling
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINELISTTAG})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTTAG))
		// inspect using the tag we used for pulling
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINELISTTAG})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTDIGEST))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINEARM64DIGEST))
		// inspect using the digest of the list
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINELISTDIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTTAG))
		// inspect using the digest of the list
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINELISTDIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTDIGEST))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINEARM64DIGEST))
		// inspect using the digest of the arch-specific image's manifest
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINEARM64DIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTTAG))
		// inspect using the digest of the arch-specific image's manifest
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINEARM64DIGEST})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTDIGEST))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINEARM64DIGEST))
		// inspect using the image ID
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoTags}}", ALPINEARM64ID})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTTAG))
		// inspect using the image ID
		session = podmanTest.PodmanNoCache([]string{"inspect", "--format", "{{.RepoDigests}}", ALPINEARM64ID})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINELISTDIGEST))
		Expect(string(session.Out.Contents())).To(ContainSubstring(ALPINEARM64DIGEST))
		// remove using the tag
		session = podmanTest.PodmanNoCache([]string{"rmi", ALPINELISTTAG})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull bogus image", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "umohnani/get-started"})
		session.WaitWithDefaultTimeout()
		Expect(session).To(ExitWithError())
	})

	It("podman pull from docker-archive", func() {
		SkipIfRemote("podman-remote does not support pulling from docker-archive")

		podmanTest.RestoreArtifact(ALPINE)
		tarfn := filepath.Join(podmanTest.TempDir, "alp.tar")
		session := podmanTest.PodmanNoCache([]string{"save", "-o", tarfn, "alpine"})
		session.WaitWithDefaultTimeout()

		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"rmi", "alpine"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"pull", fmt.Sprintf("docker-archive:%s", tarfn)})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"rmi", "alpine"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		// Pulling a multi-image archive without further specifying
		// which image _must_ error out. Pulling is restricted to one
		// image.
		session = podmanTest.PodmanNoCache([]string{"pull", fmt.Sprintf("docker-archive:./testdata/image/docker-two-images.tar.xz")})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(125))
		expectedError := "Unexpected tar manifest.json: expected 1 item, got 2"
		found, _ := session.ErrorGrepString(expectedError)
		Expect(found).To(Equal(true))

		// Now pull _one_ image from a multi-image archive via the name
		// and index syntax.
		session = podmanTest.PodmanNoCache([]string{"pull", fmt.Sprintf("docker-archive:./testdata/image/docker-two-images.tar.xz:@0")})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.PodmanNoCache([]string{"pull", fmt.Sprintf("docker-archive:./testdata/image/docker-two-images.tar.xz:example.com/empty:latest")})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.PodmanNoCache([]string{"pull", fmt.Sprintf("docker-archive:./testdata/image/docker-two-images.tar.xz:@1")})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.PodmanNoCache([]string{"pull", fmt.Sprintf("docker-archive:./testdata/image/docker-two-images.tar.xz:example.com/empty/but:different")})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		// Now check for some errors.
		session = podmanTest.PodmanNoCache([]string{"pull", fmt.Sprintf("docker-archive:./testdata/image/docker-two-images.tar.xz:foo.com/does/not/exist:latest")})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(125))
		expectedError = "Tag \"foo.com/does/not/exist:latest\" not found"
		found, _ = session.ErrorGrepString(expectedError)
		Expect(found).To(Equal(true))

		session = podmanTest.PodmanNoCache([]string{"pull", fmt.Sprintf("docker-archive:./testdata/image/docker-two-images.tar.xz:@2")})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(125))
		expectedError = "Invalid source index @2, only 2 manifest items available"
		found, _ = session.ErrorGrepString(expectedError)
		Expect(found).To(Equal(true))
	})

	It("podman pull from oci-archive", func() {
		SkipIfRemote("podman-remote does not support pulling from oci-archive")

		podmanTest.RestoreArtifact(ALPINE)
		tarfn := filepath.Join(podmanTest.TempDir, "oci-alp.tar")
		session := podmanTest.PodmanNoCache([]string{"save", "--format", "oci-archive", "-o", tarfn, "alpine"})
		session.WaitWithDefaultTimeout()

		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"rmi", "alpine"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"pull", fmt.Sprintf("oci-archive:%s", tarfn)})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"rmi", "alpine"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull from local directory", func() {
		SkipIfRemote("podman-remote does not support pulling from local directory")

		podmanTest.RestoreArtifact(ALPINE)
		dirpath := filepath.Join(podmanTest.TempDir, "alpine")
		os.MkdirAll(dirpath, os.ModePerm)
		imgPath := fmt.Sprintf("dir:%s", dirpath)

		session := podmanTest.PodmanNoCache([]string{"push", "alpine", imgPath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"rmi", "alpine"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"pull", imgPath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"images"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(session.LineInOutputContainsTag(filepath.Join("localhost", dirpath), "latest")).To(BeTrue())
		session = podmanTest.PodmanNoCache([]string{"rmi", "alpine"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull from local OCI directory", func() {
		SkipIfRemote("podman-remote does not support pulling from OCI directory")

		podmanTest.RestoreArtifact(ALPINE)
		dirpath := filepath.Join(podmanTest.TempDir, "alpine")
		os.MkdirAll(dirpath, os.ModePerm)
		imgPath := fmt.Sprintf("oci:%s", dirpath)

		session := podmanTest.PodmanNoCache([]string{"push", "alpine", imgPath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"rmi", "alpine"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"pull", imgPath})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		session = podmanTest.PodmanNoCache([]string{"images"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(session.LineInOutputContainsTag(filepath.Join("localhost", dirpath), "latest")).To(BeTrue())
		session = podmanTest.PodmanNoCache([]string{"rmi", "alpine"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("podman pull check quiet", func() {
		podmanTest.RestoreArtifact(ALPINE)
		setup := podmanTest.PodmanNoCache([]string{"images", ALPINE, "-q", "--no-trunc"})
		setup.WaitWithDefaultTimeout()
		Expect(setup.ExitCode()).To(Equal(0))
		shortImageId := strings.Split(setup.OutputToString(), ":")[1]

		rmi := podmanTest.PodmanNoCache([]string{"rmi", ALPINE})
		rmi.WaitWithDefaultTimeout()
		Expect(rmi.ExitCode()).To(Equal(0))

		pull := podmanTest.PodmanNoCache([]string{"pull", "-q", ALPINE})
		pull.WaitWithDefaultTimeout()
		Expect(pull.ExitCode()).To(Equal(0))

		Expect(pull.OutputToString()).To(ContainSubstring(shortImageId))
	})

	It("podman pull check all tags", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "--all-tags", "k8s.gcr.io/pause"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(session.LineInOuputStartsWith("Pulled Images:")).To(BeTrue())

		session = podmanTest.PodmanNoCache([]string{"images"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(len(session.OutputToStringArray())).To(BeNumerically(">", 4))
	})

	It("podman pull from docker with nonexist --authfile", func() {
		session := podmanTest.PodmanNoCache([]string{"pull", "--authfile", "/tmp/nonexist", ALPINE})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Not(Equal(0)))
	})

	It("podman pull + inspect from unqualified-search registry", func() {
		// Regression test for #6381:
		// Make sure that `pull shortname` and `inspect shortname`
		// refer to the same image.

		// We already tested pulling, so we can save some energy and
		// just restore local artifacts and tag them.
		podmanTest.RestoreArtifact(ALPINE)
		podmanTest.RestoreArtifact(BB)

		// What we want is at least two images which have the same name
		// and are prefixed with two different unqualified-search
		// registries from ../registries.conf.
		//
		// A `podman inspect $name` must yield the one from the _first_
		// matching registry in the registries.conf.
		getID := func(image string) string {
			setup := podmanTest.PodmanNoCache([]string{"image", "inspect", image})
			setup.WaitWithDefaultTimeout()
			Expect(setup.ExitCode()).To(Equal(0))

			data := setup.InspectImageJSON() // returns []inspect.ImageData
			Expect(len(data)).To(Equal(1))
			return data[0].ID
		}

		untag := func(image string) {
			setup := podmanTest.PodmanNoCache([]string{"untag", image})
			setup.WaitWithDefaultTimeout()
			Expect(setup.ExitCode()).To(Equal(0))

			setup = podmanTest.PodmanNoCache([]string{"image", "inspect", image})
			setup.WaitWithDefaultTimeout()
			Expect(setup.ExitCode()).To(Equal(0))

			data := setup.InspectImageJSON() // returns []inspect.ImageData
			Expect(len(data)).To(Equal(1))
			Expect(len(data[0].RepoTags)).To(Equal(0))
		}

		tag := func(image, tag string) {
			setup := podmanTest.PodmanNoCache([]string{"tag", image, tag})
			setup.WaitWithDefaultTimeout()
			Expect(setup.ExitCode()).To(Equal(0))
			setup = podmanTest.PodmanNoCache([]string{"image", "exists", tag})
			setup.WaitWithDefaultTimeout()
			Expect(setup.ExitCode()).To(Equal(0))
		}

		image1 := getID(ALPINE)
		image2 := getID(BB)

		// $ head -n2 ../registries.conf
		// [registries.search]
		// registries = ['docker.io', 'quay.io', 'registry.fedoraproject.org']
		registries := []string{"docker.io", "quay.io", "registry.fedoraproject.org"}
		name := "foo/test:tag"
		tests := []struct {
			// tag1 has precedence (see list above) over tag2 when
			// doing an inspect on "test:tag".
			tag1, tag2 string
		}{
			{
				fmt.Sprintf("%s/%s", registries[0], name),
				fmt.Sprintf("%s/%s", registries[1], name),
			},
			{
				fmt.Sprintf("%s/%s", registries[0], name),
				fmt.Sprintf("%s/%s", registries[2], name),
			},
			{
				fmt.Sprintf("%s/%s", registries[1], name),
				fmt.Sprintf("%s/%s", registries[2], name),
			},
		}

		for _, t := range tests {
			// 1) untag both images
			// 2) tag them according to `t`
			// 3) make sure that an inspect of `name` returns `image1` with `tag1`
			untag(image1)
			untag(image2)
			tag(image1, t.tag1)
			tag(image2, t.tag2)

			setup := podmanTest.PodmanNoCache([]string{"image", "inspect", name})
			setup.WaitWithDefaultTimeout()
			Expect(setup.ExitCode()).To(Equal(0))

			data := setup.InspectImageJSON() // returns []inspect.ImageData
			Expect(len(data)).To(Equal(1))
			Expect(len(data[0].RepoTags)).To(Equal(1))
			Expect(data[0].RepoTags[0]).To(Equal(t.tag1))
			Expect(data[0].ID).To(Equal(image1))
		}
	})
})
