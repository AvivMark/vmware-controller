package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/cors"
)

// Predefined directory to search for VMs
var vmxDirectory = "F:/VMS"

// LoadVMXPaths scans the specified directory for all .vmx files
func LoadVMXPaths() (map[string]string, error) {
	vms := make(map[string]string)

	// Walk through the directory and find all .vmx files
	files, err := ioutil.ReadDir(vmxDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read VM directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // skip directories
		}

		// Only consider .vmx files
		if strings.HasSuffix(file.Name(), ".vmx") {
			vmName := strings.TrimSuffix(file.Name(), ".vmx")
			vmxPath := filepath.Join(vmxDirectory, file.Name())
			vms[vmName] = vmxPath
		}
	}

	return vms, nil
}

// GetVMs lists all available VMs
func GetVMs(w http.ResponseWriter, r *http.Request) {
	vms, err := LoadVMXPaths()
	if err != nil {
		http.Error(w, "Failed to load VMs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vms)
}

// CreateVM creates a new VM based on a template .vmx file
func CreateVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	// Define template path and new VM path
	templatePath := filepath.Join(vmxDirectory, "template.vmx") // Specify your template .vmx file
	newVmxPath := filepath.Join(vmxDirectory, fmt.Sprintf("%s.vmx", vmName))

	// Check if the VM already exists
	if _, err := os.Stat(newVmxPath); !os.IsNotExist(err) {
		http.Error(w, "VM already exists", http.StatusConflict)
		return
	}

	// Copy the template file to create a new VM
	input, err := ioutil.ReadFile(templatePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read template: %v", err), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(newVmxPath, input, 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create VM: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "VM '%s' created successfully at path: %s", vmName, newVmxPath)
}

// DeleteVM deletes a VM by removing its .vmx file
func DeleteVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	vmxPath := filepath.Join(vmxDirectory, fmt.Sprintf("%s.vmx", vmName))

	// Check if the VM exists
	if _, err := os.Stat(vmxPath); os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("VM '%s' not found", vmName), http.StatusNotFound)
		return
	}

	// Delete the VM file
	err := os.Remove(vmxPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete VM: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "VM '%s' deleted successfully", vmName)
}

// StartVM starts a VM using its name (server finds the .vmx path)
func StartVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	// Load available VMs
	vms, err := LoadVMXPaths()
	if err != nil {
		http.Error(w, "Failed to load VMs", http.StatusInternalServerError)
		return
	}

	// Find the VM's .vmx path
	vmxPath, exists := vms[vmName]
	if !exists {
		http.Error(w, fmt.Sprintf("VM '%s' not found", vmName), http.StatusNotFound)
		return
	}

	// Start the VM using vmrun
	cmd := exec.Command("vmrun", "start", vmxPath)
	err = cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start VM: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "VM '%s' started at path: %s", vmName, vmxPath)
}

// StopVM stops a VM using its name (server finds the .vmx path)
func StopVM(w http.ResponseWriter, r *http.Request) {
	vmName := r.URL.Query().Get("name")
	if vmName == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	// Load available VMs
	vms, err := LoadVMXPaths()
	if err != nil {
		http.Error(w, "Failed to load VMs", http.StatusInternalServerError)
		return
	}

	// Find the VM's .vmx path
	vmxPath, exists := vms[vmName]
	if !exists {
		http.Error(w, fmt.Sprintf("VM '%s' not found", vmName), http.StatusNotFound)
		return
	}

	// Stop the VM using vmrun
	cmd := exec.Command("vmrun", "stop", vmxPath)
	err = cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to stop VM: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "VM '%s' stopped at path: %s", vmName, vmxPath)
}

func main() {
	// Create a new mux for your routes
	mux := http.NewServeMux()

	mux.HandleFunc("/start", StartVM)
	mux.HandleFunc("/stop", StopVM)
	mux.HandleFunc("/create", CreateVM)
	mux.HandleFunc("/delete", DeleteVM)
	mux.HandleFunc("/vms", GetVMs)

	// Wrap the mux with the CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow React app
		AllowedMethods:   []string{"GET", "POST", "DELETE"}, // Allow necessary HTTP methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Start the server with CORS
	fmt.Println("Go server started on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(mux)))
}
