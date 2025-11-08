import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";

interface FormatSelectorProps {
  selectedFormat: string;
  onFormatChange: (format: string) => void;
}

const formats = [
  { value: "mp4", label: "MP4" },
  { value: "avi", label: "AVI" },
  { value: "mov", label: "MOV" },
  { value: "webm", label: "WebM" },
  { value: "mkv", label: "MKV" },
];

export function FormatSelector({
  selectedFormat,
  onFormatChange,
}: FormatSelectorProps) {
  return (
    <div className="space-y-3">
      <Label htmlFor="format-select" className="text-base font-semibold">
        Output Format
      </Label>
      <Select value={selectedFormat} onValueChange={onFormatChange}>
        <SelectTrigger
          id="format-select"
          className="w-full"
          data-testid="select-format"
        >
          <SelectValue placeholder="Select output format" />
        </SelectTrigger>
        <SelectContent>
          {formats.map((format) => (
            <SelectItem
              key={format.value}
              value={format.value}
              data-testid={`option-format-${format.value}`}
            >
              {format.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
