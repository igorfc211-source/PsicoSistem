<script lang="ts">
	import type { LearnerReport } from '../../domain/types';
	import { RichTextEditor } from '$lib/shared/components';
	import { formatDateTime } from '$lib/shared/formatters';

	let {
		reports,
		onAddReport,
		onRemoveReport
	} = $props<{
		reports: LearnerReport[];
		onAddReport: (text: string) => void;
		onRemoveReport: (id: string) => void;
	}>();

	let reportText = $state('');

	// Controla a criacao de relatorios para manter data/hora no fluxo central da pagina.
	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		const text = reportText.replace(/<[^>]*>/g, '').trim();
		if (!text) return;

		onAddReport(reportText);
		reportText = '';
	}
</script>

<section class="tab-panel">
	<!-- Editor de novo relatorio: recebe HTML formatado do componente compartilhado. -->
	<form class="report-form" onsubmit={handleSubmit}>
		<RichTextEditor
			value={reportText}
			placeholder="Escreva um relatorio com formatacao."
			onChange={(value) => (reportText = value)}
		/>
		<button type="submit" class="primary-button">Salvar relatorio</button>
	</form>

	<!-- Historico de relatorios: mais recentes aparecem primeiro pelo fluxo da pagina. -->
	<div class="report-list">
		{#each reports as report}
			<article class="card">
				<!-- Cabecalho do registro: data de criacao e acao de exclusao. -->
				<div>
					<strong>{formatDateTime(report.createdAt)}</strong>
					<button type="button" class="danger-button" onclick={() => onRemoveReport(report.id)}>
						Excluir
					</button>
				</div>

				<!-- Conteudo formatado salvo pelo editor rico. -->
				<div class="report-content">{@html report.text}</div>
				<small>
					Feito em {formatDateTime(report.createdAt)} - atualizado em
					{formatDateTime(report.updatedAt)}
				</small>
			</article>
		{:else}
			<p class="empty-state">Nenhum relatorio salvo.</p>
		{/each}
	</div>
</section>
