<script lang="ts">
	import { MAX_DOCUMENT_BYTES, type LearnerDocument } from '$lib/modules/learners';
	import { formatDateTime, formatFileSize } from '$lib/shared/formatters';

	let {
		documents,
		isUploading,
		onUpload,
		onDownload,
		onRemove
	} = $props<{
		documents: LearnerDocument[];
		isUploading: boolean;
		onUpload: (event: Event) => void | Promise<void>;
		onDownload: (document: LearnerDocument) => void | Promise<void>;
		onRemove: (document: LearnerDocument) => void | Promise<void>;
	}>();
</script>

<section class="tab-panel">
	<div class="document-toolbar">
		<span>Limite por documento: {formatFileSize(MAX_DOCUMENT_BYTES)}</span>
		<label class="upload-button">
			<input type="file" multiple onchange={onUpload} disabled={isUploading} />
			<span>{isUploading ? 'Enviando...' : 'Adicionar documentos'}</span>
		</label>
	</div>

	<div class="document-list">
		{#each documents as document}
			<div class="document-row card">
				<div>
					<strong>{document.name}</strong>
					<span>
						{formatFileSize(document.storedSize)}
						{document.compressed ? ' - comprimido' : ''}
						- {formatDateTime(document.createdAt)}
					</span>
				</div>
				<div>
					<button type="button" class="secondary-button" onclick={() => onDownload(document)}>
						Baixar
					</button>
					<button type="button" class="danger-button" onclick={() => onRemove(document)}>
						Excluir
					</button>
				</div>
			</div>
		{:else}
			<p class="empty-state">Nenhum documento armazenado.</p>
		{/each}
	</div>
</section>
