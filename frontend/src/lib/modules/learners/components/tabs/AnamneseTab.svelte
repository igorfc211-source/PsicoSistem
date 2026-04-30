<script lang="ts">
	import type { LearnerDocument } from '../../domain/types';
	import { formatDateTime, formatFileSize } from '$lib/shared/formatters';

	let {
		value,
		documents,
		isUploading,
		onChange,
		onUpload,
		onDownload,
		onRemove
	} = $props<{
		value: string;
		documents: LearnerDocument[];
		isUploading: boolean;
		onChange: (value: string) => void;
		onUpload: (event: Event) => void | Promise<void>;
		onDownload: (document: LearnerDocument) => void | Promise<void>;
		onRemove: (document: LearnerDocument) => void | Promise<void>;
	}>();
</script>

<section class="tab-panel">
	<!-- Campo principal da anamnese: texto livre salvo diretamente no prontuario. -->
	<textarea
		class="large-textarea"
		value={value}
		oninput={(event) => onChange((event.currentTarget as HTMLTextAreaElement).value)}
		placeholder="Anamnese do aprendente."
	></textarea>

	<!-- Area de anexos da anamnese: arquivos ficam separados dos documentos gerais. -->
	<div class="anamnese-files card">
		<!-- Toolbar de upload: mostra estado de envio e aceita multiplos arquivos. -->
		<div class="document-toolbar">
			<span>Anexos da anamnese</span>
			<label class="upload-button">
				<input type="file" multiple onchange={onUpload} disabled={isUploading} />
				<span>{isUploading ? 'Enviando...' : 'Adicionar arquivo'}</span>
			</label>
		</div>

		<!-- Lista de arquivos anexados: permite baixar ou excluir item por item. -->
		<div class="document-list">
			{#each documents as document}
				<div class="document-row">
					<!-- Metadados do arquivo: nome, tamanho armazenado, compressao e data. -->
					<div>
						<strong>{document.name}</strong>
						<span>
							{formatFileSize(document.storedSize)}
							{document.compressed ? ' - comprimido' : ''}
							- {formatDateTime(document.createdAt)}
						</span>
					</div>

					<!-- Acoes do arquivo: download do blob local ou remocao da anamnese. -->
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
				<p class="empty-state">Nenhum arquivo anexado a anamnese.</p>
			{/each}
		</div>
	</div>
</section>
