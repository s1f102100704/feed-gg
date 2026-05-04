export type Label = {
  id: number;
  name: string;
};

export type PlayerLabelSummary = Label & {
  voteCount: number;
};

export type PlayerLabelsResponse = {
  labels: PlayerLabelSummary[];
  totalVotes: number;
};

export type PlayerLabelVoteResponse = PlayerLabelsResponse & {
  selectedLabel: PlayerLabelSummary;
};
