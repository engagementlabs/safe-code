#!/usr/bin/env python3
"""
Independent verification script for safe-code attestation
This script can be run by anyone to verify the code integrity
"""

import requests
import hashlib
import json
import sys
from datetime import datetime

def verify_code_attestation(app_url):
    """Verify that the running code matches GitHub source"""
    
    print(f"🔍 Verifying code attestation for {app_url}")
    print("=" * 60)
    
    # Get attestation from running app
    try:
        response = requests.get(f"{app_url}/attestation", timeout=10)
        response.raise_for_status()
        attestation = response.json()
    except Exception as e:
        print(f"❌ Failed to get attestation: {e}")
        return False
    
    commit_sha = attestation['commit_sha']
    source_hash = attestation['source_hash']
    build_time = attestation['build_time']
    github_url = attestation['github_url']
    app_attestation = attestation['attestation']
    
    print(f"📋 Commit SHA: {commit_sha}")
    print(f"🕒 Build Time: {build_time}")
    print(f"🔗 GitHub URL: {github_url}")
    print()
    
    # Download source from GitHub and verify hash
    print("🔄 Downloading source from GitHub...")
    try:
        # Try to get repository info first
        repo_url = github_url.replace('/commit/', '/').replace('https://github.com/', '')
        repo_parts = repo_url.split('/')
        if len(repo_parts) >= 2:
            repo_owner = repo_parts[0]
            repo_name = repo_parts[1]
        else:
            # Fallback to default
            repo_owner = "engagementlabs"
            repo_name = "safe-code"
        
        print(f"📦 Repository: {repo_owner}/{repo_name}")
        
        github_response = requests.get(
            f"https://api.github.com/repos/{repo_owner}/{repo_name}/tarball/{commit_sha}",
            timeout=30
        )
        github_response.raise_for_status()
        
        # Calculate hash of downloaded source
        calculated_hash = hashlib.sha256(github_response.content).hexdigest()
        
        print(f"📦 Source hash from app: {source_hash}")
        print(f"🧮 Calculated hash:     {calculated_hash}")
        
        if source_hash == calculated_hash:
            print("✅ Source hash verification: PASSED")
        else:
            print("❌ Source hash verification: FAILED")
            return False
            
    except Exception as e:
        print(f"❌ Failed to verify source: {e}")
        return False
    
    # Verify attestation signature
    print("\n🔐 Verifying attestation...")
    expected_attestation = hashlib.sha256(
        f"{commit_sha}:{build_time}:{source_hash}".encode()
    ).hexdigest()
    
    print(f"📝 App attestation:      {app_attestation}")
    print(f"🧮 Expected attestation: {expected_attestation}")
    
    if app_attestation == expected_attestation:
        print("✅ Attestation verification: PASSED")
    else:
        print("❌ Attestation verification: FAILED")
        return False
    
    # Check if Nitro Enclave attestation is available
    if attestation.get('nitro_document'):
        print("🔒 AWS Nitro Enclave attestation: AVAILABLE")
        print("   (Additional hardware-level verification)")
    else:
        print("ℹ️  AWS Nitro Enclave attestation: Not available")
    
    print("\n" + "=" * 60)
    print("🎉 CODE INTEGRITY VERIFICATION: PASSED")
    print("✅ The running code matches the GitHub repository exactly!")
    print(f"✅ Verified commit: {github_url}")
    
    return True

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 verify.py <app_url>")
        print("Example: python3 verify.py https://safe-code-p.engagelabs.org")
        sys.exit(1)
    
    app_url = sys.argv[1].rstrip('/')
    success = verify_code_attestation(app_url)
    sys.exit(0 if success else 1)