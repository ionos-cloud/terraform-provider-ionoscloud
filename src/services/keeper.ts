import * as fs from "fs";
import * as path from "path";

/**
 * Keeper — saves and restores declared preserved files across destructive
 * generation steps (e.g. `rm -rf docs/`).
 *
 * Files listed in `preservedFiles` (from setup.yaml) are copied to a temp
 * location before generation and restored afterwards. This class fails
 * loudly when a declared file cannot be found during save() or is missing
 * after restore(), instead of silently skipping.
 */

export interface KeeperOptions {
  /** Absolute path to the repository root. */
  repoRoot: string;
  /** Temp directory for backups. Defaults to os tmpdir + keeper-<pid>. */
  backupDir?: string;
}

export class Keeper {
  private repoRoot: string;
  private backupDir: string;
  private savedFiles: string[] = [];

  constructor(opts: KeeperOptions) {
    this.repoRoot = opts.repoRoot;
    this.backupDir =
      opts.backupDir ?? path.join(require("os").tmpdir(), `keeper-${process.pid}`);
  }

  /**
   * Save all declared preserved files to the backup directory.
   *
   * @param preservedFiles - relative paths (from repo root) to preserve
   * @throws Error if any declared file does not exist on disk
   */
  save(preservedFiles: string[]): void {
    if (!preservedFiles || preservedFiles.length === 0) {
      console.log("[Keeper] No preserved files declared — nothing to save");
      return;
    }

    // Create backup directory
    fs.mkdirSync(this.backupDir, { recursive: true });

    const missing: string[] = [];

    for (const relPath of preservedFiles) {
      const srcPath = path.resolve(this.repoRoot, relPath);

      if (!fs.existsSync(srcPath)) {
        missing.push(relPath);
        continue;
      }

      const destPath = path.join(this.backupDir, relPath);
      fs.mkdirSync(path.dirname(destPath), { recursive: true });
      fs.copyFileSync(srcPath, destPath);
      this.savedFiles.push(relPath);
      console.log(`[Keeper] Saved: ${relPath}`);
    }

    if (missing.length > 0) {
      const msg = `[Keeper] ERROR: The following declared preservedFiles were not found during save():\n` +
        missing.map((f) => `  - ${f}`).join("\n") +
        `\nThis may indicate a misconfiguration in setup.yaml or a first-time release.` +
        `\nDeclared files should exist before generation runs.`;
      console.error(msg);
      throw new Error(msg);
    }
  }

  /**
   * Restore all previously saved files back to the repository.
   *
   * @throws Error if any file that was successfully saved cannot be restored
   */
  restore(): void {
    if (this.savedFiles.length === 0) {
      console.log("[Keeper] No files were saved — nothing to restore");
      return;
    }

    const failedRestores: string[] = [];

    for (const relPath of this.savedFiles) {
      const srcPath = path.join(this.backupDir, relPath);
      const destPath = path.resolve(this.repoRoot, relPath);

      if (!fs.existsSync(srcPath)) {
        console.error(`[Keeper] ERROR: Backup file missing: ${srcPath}`);
        failedRestores.push(relPath);
        continue;
      }

      fs.mkdirSync(path.dirname(destPath), { recursive: true });
      fs.copyFileSync(srcPath, destPath);
      console.log(`[Keeper] Restored: ${relPath}`);

      // Verify the file actually landed
      if (!fs.existsSync(destPath)) {
        console.error(`[Keeper] ERROR: File not found after restore: ${destPath}`);
        failedRestores.push(relPath);
      }
    }

    // Clean up backup directory
    this.cleanup();

    if (failedRestores.length > 0) {
      const msg = `[Keeper] ERROR: The following files could not be restored:\n` +
        failedRestores.map((f) => `  - ${f}`).join("\n");
      console.error(msg);
      throw new Error(msg);
    }

    console.log(`[Keeper] All ${this.savedFiles.length} file(s) restored successfully`);
  }

  /** Remove the temporary backup directory. */
  private cleanup(): void {
    try {
      fs.rmSync(this.backupDir, { recursive: true, force: true });
    } catch {
      console.warn(`[Keeper] Warning: could not clean up backup dir: ${this.backupDir}`);
    }
  }
}
