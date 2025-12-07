#!/usr/bin/env python3
"""
对比原项目 (skyblessings) 和 新项目 (skyblessings-api) 的 PNG 图片
依赖库: pip install pillow
"""

from PIL import Image
import os

def get_png_info(filepath):
    """获取 PNG 文件的关键信息"""
    try:
        img = Image.open(filepath)
        if img.mode not in ('RGBA', 'LA', 'PA'):
            return None
        
        if img.mode == 'RGBA':
            alpha = img.split()[3]
        elif img.mode == 'LA':
            alpha = img.split()[1]
        else:
            alpha = img.split()[1]
        
        alpha_values = list(alpha.getdata())
        
        # 判断是否为纯遮罩（只有 0 和 255）
        unique_alphas = set(alpha_values)
        is_pure_mask = unique_alphas <= {0, 255}
        
        transparent = sum(1 for a in alpha_values if a == 0)
        opaque = sum(1 for a in alpha_values if a == 255)
        
        return {
            'mode': img.mode,
            'size': (img.width, img.height),
            'is_pure_mask': is_pure_mask,
            'unique_alphas': len(unique_alphas),
            'transparent_ratio': transparent / len(alpha_values) * 100,
            'opaque_ratio': opaque / len(alpha_values) * 100,
        }
    except:
        return None


def main():
    original_dir = r"G:\GGames\Minecraft\shuyeyun\qq-bot\xingwo\skyblessings\starimg"
    new_dir = r"G:\GGames\Minecraft\shuyeyun\qq-bot\xingwo\skyblessings-api\assets\image"
    
    print("🔍 原项目 vs 新项目 PNG 对比")
    print("="*80)
    
    # 获取所有要对比的图片
    common_names = ['background.png', 'background0.png', 'background1.png', 
                    'background2.png', 'background3.png', 'background5.png']
    
    for name in common_names:
        orig_path = os.path.join(original_dir, name)
        new_path = os.path.join(new_dir, name)
        
        print(f"\n📄 {name}")
        print("-" * 80)
        
        orig_info = get_png_info(orig_path) if os.path.exists(orig_path) else None
        new_info = get_png_info(new_path) if os.path.exists(new_path) else None
        
        if orig_info:
            print(f"  原项目: ", end="")
            if orig_info['is_pure_mask']:
                print(f"✓ 纯遮罩 (只有 0 和 255)")
            else:
                print(f"✗ 装饰性背景 ({orig_info['unique_alphas']} 种 Alpha 值)")
            print(f"           透明: {orig_info['transparent_ratio']:.1f}% | 不透明: {orig_info['opaque_ratio']:.1f}%")
        else:
            print(f"  原项目: ❌ 不存在或读取失败")
        
        if new_info:
            print(f"  新项目: ", end="")
            if new_info['is_pure_mask']:
                print(f"✓ 纯遮罩 (只有 0 和 255)")
            else:
                print(f"✗ 装饰性背景 ({new_info['unique_alphas']} 种 Alpha 值)")
            print(f"           透明: {new_info['transparent_ratio']:.1f}% | 不透明: {new_info['opaque_ratio']:.1f}%")
        else:
            print(f"  新项目: ❌ 不存在或读取失败")
    
    print("\n" + "="*80)
    print("📌 结论:")
    print("  如果 background0/1/2/3.png 都是装饰性背景（不是纯遮罩），")
    print("  那么在 Go 中，我们需要 for 每个 background*.png 都使用遮罩处理！")


if __name__ == "__main__":
    main()
