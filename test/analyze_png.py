#!/usr/bin/env python3
"""
分析 PNG 图片的结构，包括透明通道、颜色模式等
依赖库: pip install pillow
"""

from PIL import Image
import os
import sys

def analyze_png(filepath):
    """分析一个 PNG 文件的属性"""
    if not os.path.exists(filepath):
        print(f"❌ 文件不存在: {filepath}")
        return
    
    try:
        img = Image.open(filepath)
        filename = os.path.basename(filepath)
        
        print(f"\n📄 {filename}")
        print(f"  {'='*60}")
        print(f"  尺寸: {img.width} × {img.height}")
        print(f"  模式: {img.mode}")  # RGB, RGBA, L, LA, P, PA 等
        print(f"  格式: {img.format}")
        
        # 分析 Alpha 通道
        if img.mode in ('RGBA', 'LA', 'PA'):
            print(f"  ✓ 有透明通道 (Alpha Channel)")
            
            # 提取 Alpha 通道
            if img.mode == 'RGBA':
                alpha = img.split()[3]
            elif img.mode == 'LA':
                alpha = img.split()[1]
            elif img.mode == 'PA':
                alpha = img.split()[1]
            
            # 分析 Alpha 值的分布
            alpha_values = list(alpha.getdata())
            alpha_min = min(alpha_values)
            alpha_max = max(alpha_values)
            alpha_avg = sum(alpha_values) / len(alpha_values)
            
            # 计算透明、半透明、不透明的像素数
            transparent_count = sum(1 for a in alpha_values if a == 0)
            opaque_count = sum(1 for a in alpha_values if a == 255)
            semi_count = len(alpha_values) - transparent_count - opaque_count
            
            print(f"  Alpha 值范围: {alpha_min} - {alpha_max}")
            print(f"  Alpha 平均值: {alpha_avg:.1f}")
            print(f"  透明像素 (α=0): {transparent_count} ({transparent_count/len(alpha_values)*100:.1f}%)")
            print(f"  半透明像素 (0<α<255): {semi_count} ({semi_count/len(alpha_values)*100:.1f}%)")
            print(f"  不透明像素 (α=255): {opaque_count} ({opaque_count/len(alpha_values)*100:.1f}%)")
            
        else:
            print(f"  ✗ 没有透明通道")
        
        # 如果有调色板，显示调色板信息
        if img.mode in ('P', 'PA'):
            palette = img.getpalette()
            if palette:
                print(f"  调色板颜色数: {len(palette) // 3 if img.mode == 'P' else len(palette) // 4}")
        
    except Exception as e:
        print(f"❌ 分析失败: {e}")


def main():
    assets_dir = r"G:\GGames\Minecraft\shuyeyun\qq-bot\xingwo\skyblessings-api\assets\image"
    
    print("🔍 PNG 图片结构分析报告")
    print("="*60)
    
    # 分析所有 PNG 文件
    png_files = sorted([f for f in os.listdir(assets_dir) if f.endswith('.png')])
    
    for filename in png_files:
        filepath = os.path.join(assets_dir, filename)
        analyze_png(filepath)
    
    print("\n" + "="*60)
    print("📌 分析完成！")
    print("\n关键发现:")
    print("  • RGBA 模式 + 高 Alpha 变化 = 装饰性背景（不是纯遮罩）")
    print("  • RGBA 模式 + 低 Alpha 变化 = 纯遮罩（只有 0 或 255）")


if __name__ == "__main__":
    main()
