// Extrai iniciais consistentes para avatares sem depender de imagem de perfil.
export function getInitials(value?: string | null) {
	if (!value) return 'PS';

	return value
		.split(' ')
		.filter(Boolean)
		.slice(0, 2)
		.map((part) => part[0]?.toUpperCase())
		.join('');
}

// Formata data e hora de registros clinicos para leitura rapida.
export function formatDateTime(value: string) {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return value;

	return `${formatDateForNBR(date)} ${formatTimeForNBR(date)}`;
}

// Usa meio-dia local para evitar que datas puras virem o dia anterior por fuso horario.
export function formatLongDate(value: string) {
	const date = new Date(`${value}T12:00:00`);
	if (Number.isNaN(date.getTime())) return value;

	return new Intl.DateTimeFormat('pt-BR', {
		weekday: 'long',
		day: '2-digit',
		month: 'long',
		year: 'numeric'
	}).format(date);
}

// Exibe o mes corrente nos calendarios compactos e completos.
export function formatMonth(value: Date) {
	return new Intl.DateTimeFormat('pt-BR', {
		month: 'long',
		year: 'numeric'
	}).format(value);
}

// Representa datas no formato estendido da NBR 5892:2019, equivalente a ISO 8601.
export function formatDateForNBR(value: string | Date) {
	const date =
		typeof value === 'string' ? new Date(value.includes('T') ? value : `${value}T12:00:00`) : value;
	if (Number.isNaN(date.getTime())) return String(value);

	return [
		date.getFullYear(),
		String(date.getMonth() + 1).padStart(2, '0'),
		String(date.getDate()).padStart(2, '0')
	].join('-');
}

// Representa horarios em HH:mm, mantendo leitura brasileira e aderencia ao padrao NBR.
export function formatTimeForNBR(value: string | Date) {
	if (typeof value === 'string') return value.slice(0, 5);

	return [
		String(value.getHours()).padStart(2, '0'),
		String(value.getMinutes()).padStart(2, '0')
	].join(':');
}

// Mantem tamanhos de arquivo legiveis na tela de documentos.
export function formatFileSize(value: number) {
	if (value >= 1024 * 1024) return `${(value / 1024 / 1024).toFixed(1)} MB`;
	return `${Math.max(1, Math.round(value / 1024))} KB`;
}
