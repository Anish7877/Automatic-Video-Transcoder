interface FormatSelectorProps {
  selectedFormat: string;
  onFormatChange: (format: string) => void;
}

const formats = [
  { value: 'mp4', label: 'MP4', description: 'Universal compatibility' },
  { value: 'avi', label: 'AVI', description: 'High quality, large files' },
  { value: 'mov', label: 'MOV', description: 'Apple QuickTime' },
  { value: 'mkv', label: 'MKV', description: 'Open source, feature-rich' },
  { value: 'webm', label: 'WebM', description: 'Web optimized' },
  { value: 'flv', label: 'FLV', description: 'Flash video' },
  { value: 'wmv', label: 'WMV', description: 'Windows Media' },
  { value: 'mpeg', label: 'MPEG', description: 'Standard video format' },
];

function FormatSelector({ selectedFormat, onFormatChange }: FormatSelectorProps) {
  return (
    <div>
      <label className="block text-sm font-medium text-slate-700 mb-3">
        Select output format
      </label>
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
        {formats.map((format) => (
          <button
            key={format.value}
            onClick={() => onFormatChange(format.value)}
            className={`p-4 rounded-lg border-2 transition-all duration-200 text-left ${
              selectedFormat === format.value
                ? 'border-slate-900 bg-slate-50'
                : 'border-slate-200 hover:border-slate-300 bg-white'
            }`}
          >
            <div className="font-semibold text-slate-900 mb-1">
              {format.label}
            </div>
            <div className="text-xs text-slate-500 font-light">
              {format.description}
            </div>
          </button>
        ))}
      </div>
    </div>
  );
}

export default FormatSelector;
