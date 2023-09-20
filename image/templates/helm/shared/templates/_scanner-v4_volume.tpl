{{/*
  srox.scannerV4Volume $

  Configures and initializes Scanner v4 persistence.
   */}}
{{ define "srox.scannerV4Volume" }}
{{ $ := . }}
{{ $_ := set $ "_rox" $._rox }}
{{ $scannerV4DBCfg := $._rox.scannerV4.db }}

{{/*
    Scanner v4 DB Volume config setup.
  */}}
{{ $scannerV4DBVolumeCfg := dict }}
{{ if $scannerV4DBCfg.persistence.none }}
  {{ include "srox.warn" (list $ "You have selected no persistence backend. Every deletion of the StackRox Scanner v4 DB pod will cause you to lose all your data. This is STRONGLY recommended against.") }}
  {{ $_ := set $scannerV4DBVolumeCfg "emptyDir" dict }}
{{ end }}
{{ if $scannerV4DBCfg.persistence.hostPath }}
  {{ if not $scannerV4DBCfg.nodeSelector }}
    {{ include "srox.warn" (list $ "You have selected host path persistence, but not specified a node selector. This is unlikely to work reliably.") }}
  {{ end }}
  {{ $_ := set $scannerV4DBVolumeCfg "hostPath" (dict "path" $scannerV4DBCfg.persistence.hostPath) }}
{{ end }}
{{/* Configure PVC if any of the settings in `scannerV4.db.persistence.persistentVolumeClaim` is provided
     or no other persistence backend has been configured yet. */}}
{{ if or (not (deepEqual $._rox._configShapeScannerV4.scannerV4.db.persistence.persistentVolumeClaim $scannerV4DBCfg.persistence.persistentVolumeClaim)) (not $scannerV4DBVolumeCfg) }}
  {{ $scannerV4DBPVCCfg := $scannerV4DBCfg.persistence.persistentVolumeClaim }}
  {{ $_ := include "srox.mergeInto" (list $scannerV4DBPVCCfg $._rox._defaults.scannerV4DBPVCDefaults (dict "createClaim" .Release.IsInstall)) }}
  {{ $_ = set $scannerV4DBVolumeCfg "persistentVolumeClaim" (dict "claimName" $scannerV4DBPVCCfg.claimName) }}
  {{ if $scannerV4DBPVCCfg.createClaim }}
    {{ $_ = set $scannerV4DBCfg.persistence "_pvcCfg" $scannerV4DBPVCCfg }}
  {{ end }}
{{ end }}

{{ $allPersistenceMethods := keys $scannerV4DBVolumeCfg | sortAlpha }}
{{ if ne (len $allPersistenceMethods) 1 }}
  {{ include "srox.fail" (printf "Invalid or no persistence configurations for scanner-v4-db: [%s]" (join "," $allPersistenceMethods)) }}
{{ end }}
{{ $_ = set $scannerV4DBCfg.persistence "_volumeCfg" $scannerV4DBVolumeCfg }}
{{ end }}
