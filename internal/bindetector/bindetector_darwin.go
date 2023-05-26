//go:build darwin
// +build darwin

package bindetector

func isBinary(filepath string) bool {
	// regex check the output of /usr/bin/file and see if binary is within the string...

	// os.Exec to /usr/bin/file
	// _, err := os.Exec()

	// parse regex path or string comparison to check if string contains binary
	return false
}

// da python
// def is_binary(file_path):
//     """
//   is_binary(file_path)

//   Uses 'file' system cmd to determine if a file is binary.
//   """

//     bin_regex = re.compile(r"binary")
//     # Command paths per platform
//     if platform.system() == "Windows":
//         file_cmd = __win_find_file_exe()
//     else:
//         file_cmd = "/usr/bin/file"

//     file_output = run([file_cmd, "--mime-encoding", "-b", file_path])["stdout"]
//     if bin_regex.search(file_output) is not None:
//         return True

//     return False
