import subprocess
import re
import matplotlib.pyplot as plt
import os

def run_stats_script():
    """Runs filter_stats.py and captures its output."""
    script_path = os.path.join(os.path.dirname(__file__), 'filter_stats.py')
    if not os.path.exists(script_path):
        print(f"Error: {script_path} not found.")
        return None
    try:
        result = subprocess.run(['python3', script_path], capture_output=True, text=True, check=True, encoding='utf-8')
        return result.stdout
    except FileNotFoundError:
        print(f"Error: 'python3' interpreter not found or not in PATH. Make sure Python 3 is installed and accessible.")
        return None
    except subprocess.CalledProcessError as e:
        print(f"Error running {script_path}: {e}")
        print(f"Stderr: {e.stderr}")
        return None
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        return None

def parse_tld_stats(output):
    """Parses the TLD distribution from the script output."""
    tld_counts = {}
    in_tld_section = False
    # Regex to match the TLD count lines like "  .com: 123"
    tld_line_re = re.compile(r"^\s+\.([^:]+):\s+(\d+)$")

    lines = output.strip().split('\n')
    for line in lines:
        if line.strip() == "Uzantıya göre dağılım:":
            in_tld_section = True
            continue
        if in_tld_section:
            match = tld_line_re.match(line)
            if match:
                tld = match.group(1)
                count = int(match.group(2))
                tld_counts[tld] = count
            # Stop parsing if an empty line or non-matching line is encountered after the section starts
            elif line.strip() == "":
                 break
            # Add any other potential break conditions if the output format changes
    return tld_counts

def parse_active_inactive_stats(output):
    """Parses the active/inactive domain counts from the script output."""
    stats = {}
    active_re = re.compile(r"^Aktif domain:\s+(\d+)$")
    passive_re = re.compile(r"^Pasif domain:\s+(\d+)$")
    lines = output.strip().split('\n')
    for line in lines:
        active_match = active_re.match(line)
        passive_match = passive_re.match(line)
        if active_match:
            stats['Aktif'] = int(active_match.group(1))
        elif passive_match:
            stats['Pasif'] = int(passive_match.group(1))
        # Optimization: Stop after finding both stats if they appear early
        if 'Aktif' in stats and 'Pasif' in stats:
            break
    return stats

def plot_tld_distribution(tld_counts, output_filename='tld_distribution.png'):
    """Generates and saves a bar chart of TLD distribution."""
    if not tld_counts:
        print("No TLD data to plot.")
        return

    # Sort TLDs by count for better visualization
    sorted_tlds = sorted(tld_counts.items(), key=lambda item: item[1], reverse=True)
    tlds = [item[0] for item in sorted_tlds]
    counts = [item[1] for item in sorted_tlds]

    plt.figure(figsize=(12, 8)) # Adjust figure size as needed
    bars = plt.bar(tlds, counts, color='skyblue')
    plt.xlabel('Top-Level Domain (TLD)')
    plt.ylabel('Number of Domains')
    plt.title('Domain Distribution by TLD in Filter List')
    plt.xticks(rotation=45, ha='right') # Rotate labels for better readability
    plt.tight_layout() # Adjust layout to prevent labels overlapping

    # Add counts on top of bars
    for bar in bars:
        yval = bar.get_height()
        plt.text(bar.get_x() + bar.get_width()/2.0, yval, int(yval), va='bottom', ha='center') # va: vertical alignment

    try:
        plt.savefig(output_filename)
        print(f"TLD distribution plot saved to {output_filename}")
    except Exception as e:
        print(f"Error saving TLD plot: {e}")
    plt.close() # Close the figure to free memory

def plot_active_inactive_distribution(stats, output_filename='active_inactive_distribution.png'):
    """Generates and saves a pie chart of active vs inactive domains."""
    if not stats or 'Aktif' not in stats or 'Pasif' not in stats:
        print("No Active/Inactive data to plot or data is incomplete.")
        return

    labels = ['Aktif', 'Pasif']
    sizes = [stats.get('Aktif', 0), stats.get('Pasif', 0)]
    colors = ['lightgreen', 'lightcoral']
    explode = (0.1, 0) # explode the 'Aktif' slice slightly

    # Avoid plotting if total is zero
    if sum(sizes) == 0:
        print("Total number of active and passive domains is zero. Skipping plot.")
        return

    plt.figure(figsize=(8, 6))
    plt.pie(sizes, explode=explode, labels=labels, colors=colors, autopct='%1.1f%%',
            shadow=True, startangle=90)
    plt.axis('equal') # Equal aspect ratio ensures that pie is drawn as a circle.
    plt.title('Active vs. Inactive Domain Status')

    try:
        plt.savefig(output_filename)
        print(f"Active/Inactive status plot saved to {output_filename}")
    except Exception as e:
        print(f"Error saving active/inactive plot: {e}")
    plt.close() # Close the figure

def main():
    output = run_stats_script()
    if output:
        script_dir = os.path.dirname(__file__)
        project_root = os.path.join(script_dir, '..') # Get the parent directory

        # --- TLD Plot ---
        tld_stats = parse_tld_stats(output)
        if tld_stats:
            tld_plot_filename = os.path.join(project_root, 'tld_distribution.png')
            plot_tld_distribution(tld_stats, tld_plot_filename)
        else:
            print("Could not parse TLD statistics from the output.")

        # --- Active/Inactive Plot ---
        active_inactive_stats = parse_active_inactive_stats(output)
        if active_inactive_stats:
            active_inactive_plot_filename = os.path.join(project_root, 'active_inactive_distribution.png')
            plot_active_inactive_distribution(active_inactive_stats, active_inactive_plot_filename)
        else:
            print("Could not parse Active/Inactive statistics from the output.")

if __name__ == '__main__':
    try:
        import matplotlib
        matplotlib.use('Agg') # Use Agg backend for non-interactive plotting
        import matplotlib.pyplot as plt
    except ImportError:
        print("Error: matplotlib is not installed. Please install it using 'pip install matplotlib'")
    else:
        main() 