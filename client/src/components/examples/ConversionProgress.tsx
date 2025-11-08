import { ConversionProgress } from "../ConversionProgress";

export default function ConversionProgressExample() {
  return (
    <div className="space-y-8">
      <div>
        <p className="text-sm text-muted-foreground mb-4">Converting (45%):</p>
        <ConversionProgress progress={45} isConverting={true} isComplete={false} />
      </div>
      <div>
        <p className="text-sm text-muted-foreground mb-4">Complete:</p>
        <ConversionProgress progress={100} isConverting={false} isComplete={true} />
      </div>
    </div>
  );
}
