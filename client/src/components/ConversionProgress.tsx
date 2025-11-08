import { Progress } from "@/components/ui/progress";
import { motion } from "framer-motion";
import { Loader2, CheckCircle2 } from "lucide-react";

interface ConversionProgressProps {
  progress: number;
  isConverting: boolean;
  isComplete: boolean;
}

export function ConversionProgress({
  progress,
  isConverting,
  isComplete,
}: ConversionProgressProps) {
  if (!isConverting && !isComplete) return null;

  return (
    <motion.div
      initial={{ opacity: 0, height: 0 }}
      animate={{ opacity: 1, height: "auto" }}
      exit={{ opacity: 0, height: 0 }}
      transition={{ duration: 0.3 }}
      className="space-y-4"
    >
      {isConverting && !isComplete && (
        <>
          <div
            className="flex items-center justify-center gap-3 text-muted-foreground"
            data-testid="conversion-status"
          >
            <Loader2 className="h-5 w-5 animate-spin" />
            <span className="font-medium">Converting your video...</span>
          </div>
          <Progress value={progress} className="w-full" data-testid="progress-bar" />
          <p
            className="text-center text-sm text-muted-foreground"
            data-testid="text-progress"
          >
            {progress}% complete
          </p>
        </>
      )}
      {isComplete && (
        <motion.div
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 0.3 }}
          className="flex items-center justify-center gap-3 text-primary"
          data-testid="conversion-complete"
        >
          <CheckCircle2 className="h-6 w-6" />
          <span className="font-semibold text-lg">Conversion Complete!</span>
        </motion.div>
      )}
    </motion.div>
  );
}
